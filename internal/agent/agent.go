package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ncyellow/devops/internal/server/storage"
)

const (
	pollInterval   = time.Second * 2
	reportInterval = time.Second * 10
)

// Config содержит параметры по настройке агента
type Config struct {
	// Host строка в формате localhost:8080
	Host string
}

// Metrics текущее состояние всех метрик обновляются с интервалом pollInterval
type Metrics struct {
	PollCount   int64
	RandomValue float64
	runtime.MemStats
}

// Agent опрашивает метрики и отправляет их на сервер с интервалом reportInterval.
// Пример запуска:
//	conf := agent.Config{
//		Host: "localhost:8080",
//	}
//	collector := agent.Agent{Conf: conf}
type Agent struct {
	Conf    Config
	metrics Metrics
}

// sendToServer отправка метрик на сервер
func (collector *Agent) sendToServer() {
	//! приводим все метрики к нужным типам.
	gauges := prepareGauges(&collector.metrics)
	for name, value := range gauges {
		url := fmt.Sprintf("http://%s/update/", collector.Conf.Host)

		metric := storage.Metrics{
			ID:    name,
			MType: storage.Counter,
			Value: &value,
		}

		buf, err := json.Marshal(metric)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(buf))
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
	}

	counters := prepareCounters(&collector.metrics)
	for name, value := range counters {
		url := fmt.Sprintf("http://%s/update/", collector.Conf.Host)

		metric := storage.Metrics{
			ID:    name,
			MType: storage.Counter,
			Delta: &value,
		}

		buf, err := json.Marshal(metric)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(buf))
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
	}
}

// Run запускает цикл по обработке таймеров и ожидания сигналов от ОС
func (collector *Agent) Run() error {
	tickerPoll := time.NewTicker(pollInterval)
	tickerReport := time.NewTicker(reportInterval)

	defer tickerPoll.Stop()
	defer tickerReport.Stop()

	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	for {
		select {
		case <-tickerPoll.C:
			//! Обновляем все стандартные метрики
			//! Инкремент счетчика и новый рандом
			runtime.ReadMemStats(&collector.metrics.MemStats)
			collector.metrics.PollCount += 1

			rand.Seed(time.Now().UnixNano())
			collector.metrics.RandomValue = rand.Float64()

		case <-tickerReport.C:
			collector.sendToServer()

		case <-signalChanel:
			//! Корректный выход без ошибок по указанным сигналам
			return nil
		}
	}
}

// prepareGauges - готовит map[string]float64 с метриками gauges для отправки на сервер,
// так как класс метрики довольно жирный передает через указатель
func prepareGauges(metrics *Metrics) map[string]float64 {
	return map[string]float64{
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
}

// prepareCounters - готовит map[string]int64 с метриками counter для отправки на сервер,
// пока такая метрика одна, но для обобщения сделан сразу метод
func prepareCounters(metrics *Metrics) map[string]int64 {
	return map[string]int64{
		"PollCount": metrics.PollCount,
	}
}
