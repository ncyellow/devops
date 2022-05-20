package agent

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

const (
	pollInterval   = time.Second * 2
	reportInterval = time.Second * 10
)

type Metrics struct {
	PollCount   int64
	RandomValue float64
	runtime.MemStats
}

type Agent struct {
	metrics Metrics
}

func (collector *Agent) sendToServer() {
	//! приводим все метрики к нужным типам.
	var data = map[string]float64{
		"Alloc":         float64(collector.metrics.Alloc),
		"BuckHashSys":   float64(collector.metrics.BuckHashSys),
		"Frees":         float64(collector.metrics.Frees),
		"GCCPUFraction": collector.metrics.GCCPUFraction,
		"GCSys":         float64(collector.metrics.GCSys),
		"HeapAlloc":     float64(collector.metrics.HeapAlloc),
		"HeapIdle":      float64(collector.metrics.HeapIdle),
		"HeapInuse":     float64(collector.metrics.HeapInuse),
		"HeapObjects":   float64(collector.metrics.HeapObjects),
		"HeapReleased":  float64(collector.metrics.HeapReleased),
		"HeapSys":       float64(collector.metrics.HeapSys),
		"LastGC":        float64(collector.metrics.LastGC),
		"Lookups":       float64(collector.metrics.Lookups),
		"MCacheInuse":   float64(collector.metrics.MCacheInuse),
		"MCacheSys":     float64(collector.metrics.MCacheSys),
		"MSpanInuse":    float64(collector.metrics.MSpanInuse),
		"MSpanSys":      float64(collector.metrics.MSpanSys),
		"Mallocs":       float64(collector.metrics.Mallocs),
		"NextGC":        float64(collector.metrics.NextGC),
		"NumForcedGC":   float64(collector.metrics.NumForcedGC),
		"NumGC":         float64(collector.metrics.NumGC),
		"OtherSys":      float64(collector.metrics.OtherSys),
		"PauseTotalNs":  float64(collector.metrics.PauseTotalNs),
		"StackInuse":    float64(collector.metrics.StackInuse),
		"StackSys":      float64(collector.metrics.StackSys),
		"Sys":           float64(collector.metrics.Sys),
		"TotalAlloc":    float64(collector.metrics.TotalAlloc),
		"RandomValue":   collector.metrics.RandomValue,
	}
	for name, value := range data {
		url := fmt.Sprintf("http://%s/update/gauge/%s/%f", "127.0.0.1:8080", name, value)
		http.Post(url, "text/plain", nil)
	}

	url := fmt.Sprintf("http://%s/update/counter/%s/%d",
		"127.0.0.1:8080",
		"PollCount",
		collector.metrics.PollCount)
	http.Post(url, "text/plain", nil)
}

func (collector *Agent) Run() error {
	tickerPoll := time.NewTicker(pollInterval)
	tickerReport := time.NewTicker(reportInterval)
	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	for true {
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
	return nil
}
