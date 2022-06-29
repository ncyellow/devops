package storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/ncyellow/devops/internal/server/config"
	"github.com/ncyellow/devops/internal/server/repository"
	"github.com/rs/zerolog/log"
)

// RunSaver запускает сохранение данных repo по таймеру в файл
func RunSaver(pStore PersistentStorage, interval time.Duration) {
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

func (p *PgPersistentStorage) init() {

	// Это не clickhouse, где данные хорошо жмутся по столбцам и они не зависимы
	// в postgresql таблица просто будет больше по размеру из-за пустых столбцов.

	// Создаем нужные таблицы если их нет

	rows, err := p.conn.Query(context.Background(), `
	CREATE TABLE IF NOT EXISTS "counters"(
	"@counters" bigserial,
	"metric_name" text NOT NULL,
	"value" bigint
	)`)
	rows.Close()

	if err != nil {
		log.Warn().Msgf("Ошибка при создании таблицы counters - %#v\n", err)
	}

	rows, err = p.conn.Query(context.Background(), `
	CREATE UNIQUE INDEX IF NOT EXISTS "icounters-metric_name"
	ON "counters" USING btree
	("metric_name")
	`)
	rows.Close()

	if err != nil {
		log.Warn().Msgf("Ошибка при индекса уникальности по имени метрики counters - %#v\n", err)
	}

	rows, err = p.conn.Query(context.Background(), `
	CREATE TABLE IF NOT EXISTS "gauges"(
	"@gauges" bigserial,
	"metric_name" text NOT NULL,
	"value" double precision
	)`)
	rows.Close()

	if err != nil {
		log.Warn().Msgf("Ошибка при создании таблицы gauges - %#v\n", err)

	}

	rows, err = p.conn.Query(context.Background(), `
	CREATE UNIQUE INDEX IF NOT EXISTS "igauges-metric_name"
	ON "gauges" USING btree
	("metric_name")
	`)
	rows.Close()

	if err != nil {
		log.Warn().Msgf("Ошибка при индекса уникальности по имени метрики gauges - %#v\n", err)

	}

}

func NewPGStorage(conf *config.Config, repo repository.Repository) (PersistentStorage, error) {

	conn, err := pgx.Connect(context.Background(), conf.DatabaseConn)
	if err != nil {
		return nil, err
	}
	saver := PgPersistentStorage{conf: conf, conn: conn, repo: repo}
	saver.init()
	return &saver, nil

}

func (p *PgPersistentStorage) Close() {
	p.Save()
	p.conn.Close(context.Background())
}

func (p *PgPersistentStorage) Load() error {

	metrics := make([]repository.Metrics, 0)

	rows, err := p.conn.Query(context.Background(), `
	select "metric_name", "value" FROM "counters"
	`)

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

	rows, err = p.conn.Query(context.Background(), `
	select "metric_name", "value" FROM "gauges"
	`)

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

// Save - так как конкретных требований нет и не нужно хранить историю метрик, просто удаляем все что было
// и вставляем новые значения
func (p *PgPersistentStorage) Save() error {

	metrics := p.repo.ToMetrics()

	if len(metrics) == 0 {
		return nil
	}

	_, err := p.conn.Exec(context.Background(), `
	TRUNCATE TABLE "counters"
	`)
	if err != nil {
		log.Warn().Msgf("have some error while truncate counters - %#v\n", err)
	}

	_, err = p.conn.Exec(context.Background(), `
	TRUNCATE TABLE "gauges"
	`)
	if err != nil {
		log.Warn().Msgf("have some error while truncate gauges - %#v\n", err)
	}

	counters := make([]repository.Metrics, 0)
	gauges := make([]repository.Metrics, 0)

	for _, value := range metrics {
		switch value.MType {
		case repository.Gauge:
			gauges = append(gauges, value)
		case repository.Counter:
			counters = append(counters, value)
		}
	}

	_, err = p.conn.CopyFrom(
		context.Background(),
		pgx.Identifier{"counters"},
		[]string{"metric_name", "value"},
		pgx.CopyFromSlice(len(counters), func(i int) ([]interface{}, error) {
			return []interface{}{counters[i].ID, *counters[i].Delta}, nil
		}),
	)

	if err != nil {
		log.Warn().Msgf("have some error while copyFrom counters - %#v\n", err)
	}

	_, err = p.conn.CopyFrom(
		context.Background(),
		pgx.Identifier{"gauges"},
		[]string{"metric_name", "value"},
		pgx.CopyFromSlice(len(gauges), func(i int) ([]interface{}, error) {
			return []interface{}{gauges[i].ID, *gauges[i].Value}, nil
		}),
	)

	if err != nil {
		log.Warn().Msgf("have some error while copyFrom gauges - %#v\n", err)
	}
	return nil
}
