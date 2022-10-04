package agent

import (
	"fmt"
	"testing"

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
