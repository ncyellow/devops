package storage

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/ncyellow/devops/internal/repository"
	"github.com/ncyellow/devops/internal/server/config"
)

// RunStorageSaver запускает сохранение данных pStore по таймеру с интервалом interval
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
		pStore.Save(context.Background())
	}
}

// PgPersistentStorage структура для работы с pgsql. Реализует интерфейс PersistentStorage
// Использую связку pgx + pgxpool это дает нам thread safety pool коннектов
type PgPersistentStorage struct {
	conf *config.Config
	pool *pgxpool.Pool
	repo repository.Repository
}

// NewPgStorage конструктор хранилища на основе postgresql, явно не используется, только через фабрику
func NewPgStorage(conf *config.Config, repo repository.Repository) (*PgPersistentStorage, error) {

	pool, err := pgxpool.Connect(context.Background(), conf.DatabaseConn)
	if err != nil {
		return nil, errors.New("cant connect to pgsql")
	}

	saver := PgPersistentStorage{conf: conf, pool: pool, repo: repo}
	saver.init()
	return &saver, nil
}

func (p *PgPersistentStorage) Ping() error {
	return p.pool.Ping(context.Background())
}

func (p *PgPersistentStorage) Close() {
	p.Save(context.Background())
	p.pool.Close()
}

func (p *PgPersistentStorage) Load() error {

	metrics := make([]repository.Metrics, 0)

	//! Загружаем counter метрики
	counterRows, err := p.pool.Query(context.Background(), `select "metric_name", "value" FROM "counters"`)

	if err != nil {
		return err
	}

	defer counterRows.Close()
	for counterRows.Next() {
		var metricName string
		var delta int64
		err = counterRows.Scan(&metricName, &delta)
		if err != nil {
			return err
		}

		metrics = append(metrics, repository.Metrics{
			ID:    metricName,
			MType: repository.Counter,
			Delta: &delta,
		})
	}

	if counterRows.Err() != nil {
		return err
	}

	//! Загружаем gauge метрики
	gaugeRows, err := p.pool.Query(context.Background(), `select "metric_name", "value" FROM "gauges"`)

	if err != nil {
		return err
	}
	defer gaugeRows.Close()

	for gaugeRows.Next() {
		var metricName string
		var value float64
		err = gaugeRows.Scan(&metricName, &value)
		if err != nil {
			return err
		}

		metrics = append(metrics, repository.Metrics{
			ID:    metricName,
			MType: repository.Gauge,
			Value: &value,
		})
	}

	if gaugeRows.Err() != nil {
		return err
	}

	p.repo.FromMetrics(metrics)
	return nil
}

func (p *PgPersistentStorage) Save(ctx context.Context) error {

	metrics := p.repo.ToMetrics()
	if len(metrics) == 0 {
		return nil
	}

	tx, err := p.pool.Begin(ctx)
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

	// Я бы, конечно, так не делал. Лучше уж делать truncate и через COPY писать все разом. Но раз упражнение
	// требует prepare stmt + транзакцию то ок. Но мы так получим дикие проблемы с VACUUM. В таблице будет оч много
	// мертвых кортежей
	for _, value := range metrics {
		switch value.MType {
		case repository.Gauge:
			tag, err := tx.Exec(ctx, repository.Gauge, value.ID, *value.Value)

			if err != nil || !tag.Insert() {
				log.Info().Msgf("insert gauges failed - %s", err.Error())
				if err = tx.Rollback(ctx); err != nil {
					log.Info().Msgf("update drivers: unable to rollback - %s", err.Error())
				}
				return err
			}
		case repository.Counter:
			tag, err := tx.Exec(ctx, repository.Counter, value.ID, *value.Delta)

			if err != nil || !tag.Insert() {
				log.Info().Msgf("insert counters failed - %s", err.Error())
				if err = tx.Rollback(ctx); err != nil {
					log.Info().Msgf("update drivers: unable to rollback - %s", err.Error())
				}
				return err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.Info().Msgf("update drivers: unable to commit - %s", err.Error())
		return err
	}
	return nil
}

// init инициализация базы данных. Надо бы будет оформить через migrate, но пока так
func (p *PgPersistentStorage) init() {

	// Создаем нужные таблицы если их нет и индекс для уникальности имени метрики в таблице
	queries := []string{
		`CREATE TABLE IF NOT EXISTS "counters"("@counters" bigserial, "metric_name" text NOT NULL, "value" bigint)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS "icounters-metric_name" ON "counters" USING btree ("metric_name")`,
		`CREATE TABLE IF NOT EXISTS "gauges"("@gauges" bigserial,"metric_name" text NOT NULL,"value" double precision)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS "igauges-metric_name" ON "gauges" USING btree ("metric_name")`,
	}

	for _, query := range queries {
		_, err := p.pool.Exec(context.Background(), query)
		if err != nil {
			log.Info().Msgf("Не удалось выполнить запрос подготовки базы данных %s", query)
		}
	}
}
