package agent

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/ncyellow/devops/internal/agent/config"
	"github.com/ncyellow/devops/internal/gcfg"
	"github.com/ncyellow/devops/internal/hash"
	"github.com/ncyellow/devops/internal/server/repository"
)

// RuntimeMetrics текущее состояние всех метрик обновляются с интервалом pollInterval
type RuntimeMetrics struct {
	PollCount   int64
	RandomValue float64
	runtime.MemStats
}

// prepareGauges - готовит map[string]float64 с метриками gauges для отправки на сервер,
// так как класс метрики довольно жирный передает через указатель
func (metrics *RuntimeMetrics) prepareGauges(secretKey string) []repository.Metrics {
	gauges := map[string]float64{
		"Alloc":         float64(metrics.Alloc),
		"BuckHashSys":   float64(metrics.BuckHashSys),
		"Frees":         float64(metrics.Frees),
		"GCCPUFraction": metrics.GCCPUFraction,
		"GCSys":         float64(metrics.GCSys),
		"HeapAlloc":     float64(metrics.HeapAlloc),
		"HeapIdle":      float64(metrics.HeapIdle),
		"HeapInuse":     float64(metrics.HeapInuse),
		"HeapObjects":   float64(metrics.HeapObjects),
		"HeapReleased":  float64(metrics.HeapReleased),
		"HeapSys":       float64(metrics.HeapSys),
		"LastGC":        float64(metrics.LastGC),
		"Lookups":       float64(metrics.Lookups),
		"MCacheInuse":   float64(metrics.MCacheInuse),
		"MCacheSys":     float64(metrics.MCacheSys),
		"MSpanInuse":    float64(metrics.MSpanInuse),
		"MSpanSys":      float64(metrics.MSpanSys),
		"Mallocs":       float64(metrics.Mallocs),
		"NextGC":        float64(metrics.NextGC),
		"NumForcedGC":   float64(metrics.NumForcedGC),
		"NumGC":         float64(metrics.NumGC),
		"OtherSys":      float64(metrics.OtherSys),
		"PauseTotalNs":  float64(metrics.PauseTotalNs),
		"StackInuse":    float64(metrics.StackInuse),
		"StackSys":      float64(metrics.StackSys),
		"Sys":           float64(metrics.Sys),
		"TotalAlloc":    float64(metrics.TotalAlloc),
		"RandomValue":   metrics.RandomValue,
	}
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
func (metrics *RuntimeMetrics) prepareCounters(secretKey string) []repository.Metrics {
	counters := map[string]int64{
		"PollCount": metrics.PollCount,
	}

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

type RuntimeCollector struct {
	conf    gcfg.GeneralConfig
	metrics RuntimeMetrics
}

func (rc *RuntimeCollector) Update() {
	//! Обновляем все стандартные метрики
	//! Инкремент счетчика и новый рандом
	runtime.ReadMemStats(&rc.metrics.MemStats)
	rc.metrics.PollCount += 1

	rand.Seed(time.Now().UnixNano())
	rc.metrics.RandomValue = rand.Float64()
}

func (rc *RuntimeCollector) ToMetrics() []repository.Metrics {
	allMetrics := rc.metrics.prepareGauges(rc.conf.SecretKey)
	counters := rc.metrics.prepareCounters(rc.conf.SecretKey)
	allMetrics = append(allMetrics, counters...)
	return allMetrics
}

// RunCollector запускает цикл по обработке таймеров и ожидания сигналов от ОС
func RunCollector(ctx context.Context, conf *config.Config, in chan<- []repository.Metrics, wg *sync.WaitGroup) {

	tickerPoll := time.NewTicker(conf.PollInterval)
	defer tickerPoll.Stop()

	runtimeCol := RuntimeCollector{}

	for {
		select {
		case <-tickerPoll.C:
			//! Обновляем все стандартные метрики
			//! Инкремент счетчика и новый рандом
			runtimeCol.Update()
			fmt.Println("RunCollector")
			in <- runtimeCol.ToMetrics()
		case <-ctx.Done():
			//! Корректный выход без ошибок по указанным сигналам
			wg.Done()
			return
		}
	}
}
