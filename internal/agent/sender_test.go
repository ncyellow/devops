package agent

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/ncyellow/devops/internal/agent/config"
	"github.com/ncyellow/devops/internal/genconfig"
	"github.com/ncyellow/devops/internal/repository"
)

// BenchmarkSendMetricsBatch бенчмарк на отправку метрик на сервис пачкой
func BenchmarkSendMetricsBatch(b *testing.B) {
	b.StopTimer()

	var metrics []repository.Metrics
	// Будем слать постоянно по 50 метрик за пачку
	for i := 0; i < 30; i++ {
		// делаем так, чтобы у нас все значения метрик были разные
		val := int64(i)
		metrics = append(metrics, repository.Metrics{
			ID:    fmt.Sprintf("Metric %d", i),
			MType: repository.Counter,
			Delta: &val,
		})
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		SendMetricsBatch(metrics, "http://unknown/updates/", nil)
	}
}

// Проверяем одновременный запуск трех сендеров и корректное завершение по завершению контекста всех горутин
func TestRunSender(t *testing.T) {
	// Контекст для корректно завершения все горутин
	ctx, cancel := context.WithCancel(context.Background())

	// Канал по которому метрики откуда sender модуль будет получать метрики и отправлять на сервер
	metricChannel := make(chan []repository.Metrics, 1)

	source := RuntimeSource{}
	source.Update()

	conf := &config.Config{
		GeneralConfig: genconfig.GeneralConfig{
			Address:   "localhost:8080",
			CryptoKey: "/path/to/key.pem",
		},
		ReportInterval: genconfig.Duration{Duration: time.Second * 3},
		PollInterval:   genconfig.Duration{Duration: time.Second * 32},
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go RunSender(ctx, conf, metricChannel, &wg)

	wg.Add(1)
	go RunSender(ctx, conf, metricChannel, &wg)

	wg.Add(1)
	go RunSender(ctx, conf, metricChannel, &wg)

	// Отправляем данные для отправки.
	metricChannel <- prepareGauges(source.Gauges(), "")

	time.Sleep(time.Second * 5)
	close(metricChannel)
	cancel()
	wg.Wait()

	// Если у нас будут проблемы с кривыми закрытиями каналов или отменой контекста, тест втупую зависнет
	// и сюда мы не попадем

}
