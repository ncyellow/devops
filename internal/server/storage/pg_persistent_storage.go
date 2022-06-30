package storage

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/jackc/pgx/v4"
	"github.com/ncyellow/devops/internal/server/config"
	"github.com/ncyellow/devops/internal/server/repository"
)

// RunStorageSaver запускает сохранение данных pStore по таймеру в файл
func RunStorageSaver(pStore PersistentStorage, interval time.Duration) {
	if interval == 0 {
		//! Не нужно сбрасывать на диск если StoreInterval == 0
		return
	}

	tickerStore := time.NewTicker(interval)
	defer tickerStore.Stop()

	for {
		<-tickerStore.C
		//! сбрасываем на диск
		pStore.Save()
	}
}

type PgPersistentStorage struct {
	conf *config.Config
	conn *pgx.Conn
	repo repository.Repository
}

func NewPgStorage(conf *config.Config, repo repository.Repository) (PersistentStorage, error) {

	conn, err := pgx.Connect(context.Background(), conf.DatabaseConn)
	if err != nil {
		return nil, errors.New("cant connect to pgsql")
	}

	saver := PgPersistentStorage{conf: conf, conn: conn, repo: repo}
	saver.init()
	return &saver, nil
}

func (p *PgPersistentStorage) init() {

	// Создаем нужные таблицы если их нет и индекс для уникальности имени метрики в таблице
	queries := []string{
		`CREATE TABLE IF NOT EXISTS "counters"("@counters" bigserial, "metric_name" text NOT NULL, "value" bigint)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS "icounters-metric_name" ON "counters" USING btree ("metric_name")`,
		`CREATE TABLE IF NOT EXISTS "gauges"("@gauges" bigserial,"metric_name" text NOT NULL,"value" double precision)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS "igauges-metric_name" ON "gauges" USING btree ("metric_name")`,
	}

	for _, query := range queries {
		_, err := p.conn.Exec(context.Background(), query)
		if err != nil {
			log.Info().Msgf("Не удалось выполнить запрос %s", query)
		}
	}
}

func (p *PgPersistentStorage) Close() {
	p.Save()
	p.conn.Close(context.Background())
}

func (p *PgPersistentStorage) Load() error {

	metrics := make([]repository.Metrics, 0)

	rows, err := p.conn.Query(context.Background(), `select "metric_name", "value" FROM "counters"`)

	if err != nil {
		return err
	}

	for rows.Next() {
		var metricName string
		var delta int64
		err = rows.Scan(&metricName, &delta)
		if err != nil {
			return err
		}

		metrics = append(metrics, repository.Metrics{
			ID:    metricName,
			MType: repository.Counter,
			Delta: &delta,
		})
	}

	rows.Close()

	rows, err = p.conn.Query(context.Background(), `select "metric_name", "value" FROM "gauges"`)

	if err != nil {
		return err
	}

	for rows.Next() {
		var metricName string
		var value float64
		err = rows.Scan(&metricName, &value)
		if err != nil {
			return err
		}

		metrics = append(metrics, repository.Metrics{
			ID:    metricName,
			MType: repository.Gauge,
			Value: &value,
		})
	}

	rows.Close()

	p.repo.FromMetrics(metrics)
	return nil
}

func (p *PgPersistentStorage) Save() error {

	conn, err := pgx.Connect(context.Background(), p.conf.DatabaseConn)
	if err != nil {
		return errors.New("cant connect to pgsql")
	}
	p.conn = conn

	metrics := p.repo.ToMetrics()
	if len(metrics) == 0 {
		return nil
	}

	ctx := context.Background()
	tx, err := p.conn.Begin(ctx)
	if err != nil {
		return err
	}

	desc, err := tx.Prepare(ctx, repository.Gauge, `
	INSERT INTO "gauges"("metric_name", "value")
	VALUES ($1, $2)
	ON CONFLICT ("metric_name") 
	DO 
   	UPDATE SET value = EXCLUDED.value
	`)

	if err != nil {
		log.Info().Msgf("cant create prepare stmt - %s", desc.Name)
		return err
	}

	desc, err = tx.Prepare(ctx, repository.Counter, `
	INSERT INTO "counters"("metric_name", "value")
	VALUES ($1, $2)
	ON CONFLICT ("metric_name") 
	DO 
   	UPDATE SET value = EXCLUDED.value
	`)

	if err != nil {
		log.Info().Msgf("cant create prepare stmt - %s", desc.Name)
		return err
	}

	for _, value := range metrics {
		switch value.MType {
		case repository.Gauge:
			tag, err := tx.Exec(ctx, repository.Gauge, value.ID, *value.Value)

			if err != nil || !tag.Insert() {
				log.Info().Msgf("insert gauges failed - %s", err.Error())
				if err = tx.Rollback(ctx); err != nil {
					log.Fatal().Msgf("update drivers: unable to rollback - %s", err.Error())
				}
				return err
			}
		case repository.Counter:
			tag, err := tx.Exec(ctx, repository.Counter, value.ID, *value.Delta)

			if err != nil || !tag.Insert() {
				log.Info().Msgf("insert counters failed - %s", err.Error())
				if err = tx.Rollback(ctx); err != nil {
					log.Fatal().Msgf("update drivers: unable to rollback - %s", err.Error())
				}
				return err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.Fatal().Msgf("update drivers: unable to commit - %s", err.Error())
		return err
	}

	return nil
}
