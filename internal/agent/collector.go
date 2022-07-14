package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ncyellow/devops/internal/agent/config"
	"github.com/ncyellow/devops/internal/gcfg"
	"github.com/ncyellow/devops/internal/hash"
	"github.com/ncyellow/devops/internal/server/repository"
)

type Collector struct {
	Conf   *gcfg.GeneralConfig
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

// RunCollector запускает цикл по обработке таймеров и ожидания сигналов от ОС
func RunCollector(ctx context.Context, conf *config.Config, collector *Collector, in chan<- []repository.Metrics, wg *sync.WaitGroup) {

	tickerPoll := time.NewTicker(conf.PollInterval)
	defer tickerPoll.Stop()

	for {
		select {
		case <-tickerPoll.C:
			//! Обновляем все стандартные метрики
			//! Инкремент счетчика и новый рандом
			collector.Update()
			fmt.Println("RunCollector")

			in <- collector.ToMetrics()
		case <-ctx.Done():
			//! Корректный выход без ошибок по указанным сигналам
			wg.Done()
			return
		}
	}
}

// prepareGauges - готовит map[string]float64 с метриками gauges для отправки на сервер,
// так как класс метрики довольно жирный передает через указатель
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

// prepareCounters - готовит map[string]int64 с метриками counter для отправки на сервер,
// пока такая метрика одна, но для обобщения сделан сразу метод
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
