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

	"github.com/ncyellow/devops/internal/agent/config"
	"github.com/ncyellow/devops/internal/server/storage"
)

// RuntimeMetrics текущее состояние всех метрик обновляются с интервалом pollInterval
type RuntimeMetrics struct {
	PollCount   int64
	RandomValue float64
	runtime.MemStats
}

// Agent опрашивает метрики и отправляет их на сервер с интервалом reportInterval.
// Пример запуска:
//	conf := agent.Config{
//		Address: "localhost:8080",
//	}
//	collector := agent.Agent{Conf: conf}
type Agent struct {
	Conf    *config.Config
	metrics RuntimeMetrics
}

// sendToServer отправка метрик на сервер
func (collector *Agent) sendToServer() {
	//! приводим все метрики к нужным типам.
	url := fmt.Sprintf("http://%s/update/", collector.Conf.Address)
	SendMetrics(collector.metrics.prepareGauges(), url)
	SendMetrics(collector.metrics.prepareCounters(), url)
}

// Run запускает цикл по обработке таймеров и ожидания сигналов от ОС
func (collector *Agent) Run() error {

	tickerPoll := time.NewTicker(collector.Conf.PollInterval)
	tickerReport := time.NewTicker(collector.Conf.ReportInterval)

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
func (metrics *RuntimeMetrics) prepareGauges() []storage.Metrics {
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

	result := make([]storage.Metrics, 0, len(gauges))
	for name, value := range gauges {
		// Если пользоваться value, то все значения будут ссылаться на одну и ту же переменную - последнюю
		gaugeValue := value
		metric := storage.Metrics{
			ID:    name,
			MType: storage.Gauge,
			Value: &gaugeValue,
		}
		result = append(result, metric)
	}
	return result
}

// prepareCounters - готовит map[string]int64 с метриками counter для отправки на сервер,
// пока такая метрика одна, но для обобщения сделан сразу метод
func (metrics *RuntimeMetrics) prepareCounters() []storage.Metrics {
	counters := map[string]int64{
		"PollCount": metrics.PollCount,
	}

	result := make([]storage.Metrics, 0, len(counters))
	for name, value := range counters {
		// Если пользоваться value, то все значения будут ссылаться на одну и ту же переменную - последнюю
		counterValue := value
		metric := storage.Metrics{
			ID:    name,
			MType: storage.Counter,
			Delta: &counterValue,
		}
		result = append(result, metric)
	}
	return result
}

// SendMetrics отправляет метрики на указанный url
func SendMetrics(dataSource []storage.Metrics, url string) {
	for _, metric := range dataSource {
		buf, err := json.Marshal(metric)
		if err != nil {
			log.Fatal(err)
		}
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(buf))
		if err != nil {
			log.Println(err)
			continue
		}
		resp.Body.Close()
	}
}
