// Package agent содержит АПИ для сбору и отправки метрик от разных источников
// Пример использования
// runtimeCol := &Collector{ Conf:   collector.Conf.GeneralCfg(), Source: &RuntimeSource{}, }
// go RunCollector(ctx, collector.Conf, runtimeCol, metricChannel, &wg)
package agent

import (
	"context"
	"sync"
	"time"

	"github.com/ncyellow/devops/internal/agent/config"
	"github.com/ncyellow/devops/internal/genconfig"
	"github.com/ncyellow/devops/internal/hash"
	"github.com/ncyellow/devops/internal/repository"
)

// Collector объект для работы с метриками.
// 1. Хранит источник метрик, вызывает их обновление
// 2. Готовит метрики в формат []repository.Metrics для отправки
type Collector struct {
	Conf   *genconfig.GeneralConfig
	Source MetricSource
}

func (c *Collector) Update() {
	//! Обновляем все стандартные метрики
	//! Инкремент счетчика и новый рандом
	c.Source.Update()
}

func (c *Collector) ToMetrics() []repository.Metrics {
	allMetrics := prepareGauges(c.Source.Gauges(), c.Conf.SecretKey)
	counters := prepareCounters(c.Source.Counters(), c.Conf.SecretKey)
	allMetrics = append(allMetrics, counters...)
	return allMetrics
}

// RunCollector запускает цикл по опросу метрик и отправки их в канал in
func RunCollector(ctx context.Context, conf *config.Config, collector *Collector, in chan<- []repository.Metrics, wg *sync.WaitGroup) {

	tickerPoll := time.NewTicker(conf.PollInterval)
	defer tickerPoll.Stop()

	for {
		select {
		case <-tickerPoll.C:
			// можно конечно объединить интерфейс в одну функцию, но и так ок
			collector.Update()
			in <- collector.ToMetrics()
		case <-ctx.Done():
			//! Корректный выход без ошибок по указанным сигналам
			wg.Done()
			return
		}
	}
}

// prepareGauges - преобразование метрик Gauge в []repository.Metrics
func prepareGauges(gauges map[string]float64, secretKey string) []repository.Metrics {
	hashFunc := hash.CreateEncodeFunc(secretKey)
	result := make([]repository.Metrics, 0, len(gauges))
	for name, value := range gauges {
		// Если пользоваться value, то все значения будут ссылаться на одну и ту же переменную - последнюю
		gaugeValue := value
		metric := repository.Metrics{
			ID:    name,
			MType: repository.Gauge,
			Value: &gaugeValue,
		}
		metric.Hash = metric.CalcHash(hashFunc)
		result = append(result, metric)
	}
	return result
}

// prepareCounters - преобразование метрик Counter в []repository.Metrics
func prepareCounters(counters map[string]int64, secretKey string) []repository.Metrics {
	hashFunc := hash.CreateEncodeFunc(secretKey)
	result := make([]repository.Metrics, 0, len(counters))
	for name, value := range counters {
		// Если пользоваться value, то все значения будут ссылаться на одну и ту же переменную - последнюю
		counterValue := value
		metric := repository.Metrics{
			ID:    name,
			MType: repository.Counter,
			Delta: &counterValue,
		}
		metric.Hash = metric.CalcHash(hashFunc)
		result = append(result, metric)
	}
	return result
}
