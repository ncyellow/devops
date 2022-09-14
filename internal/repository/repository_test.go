package repository

import (
	"testing"

	"github.com/ncyellow/devops/internal/hash"
)

func BenchmarkMetrics_CalcHash(b *testing.B) {
	b.StopTimer()
	ef := hash.CreateEncodeFunc("superSecret")
	val := int64(10)
	m := Metrics{
		ID:    "testCounter",
		MType: Counter,
		Delta: &val,
	}
	b.StartTimer()

	// Проверяем скорость вычисление хеша
	for i := 0; i < b.N; i++ {
		m.CalcHash(ef)
	}
}
