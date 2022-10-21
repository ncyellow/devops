package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncyellow/devops/internal/hash"
)

func TestMetrics_CalcHash(t *testing.T) {
	//! Проверяем что хеш от не адекватного типа метрики не считается
	metric := Metrics{
		ID:    "unknownMTypeMetric",
		MType: "unknownMTypeMetric",
	}
	encodeFunc := hash.CreateEncodeFunc("superKey")
	assert.Equal(t, metric.CalcHash(encodeFunc), "")
}

//func BenchmarkMetrics_CalcHash(b *testing.B) {
//	b.StopTimer()
//	ef := hash.CreateEncodeFunc("superSecret")
//	val := int64(10)
//	m := Metrics{
//		ID:    "testCounter",
//		MType: Counter,
//		Delta: &val,
//	}
//	b.StartTimer()
//
//	// Проверяем скорость вычисление хеша
//	for i := 0; i < b.N; i++ {
//		m.CalcHash(ef)
//	}
//}
