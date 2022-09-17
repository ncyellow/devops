package agent

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestRuntimeSource проверяем что источник возвращает корректное число метрик
// counter и gauge
func TestRuntimeSource(t *testing.T) {
	source := RuntimeSource{}
	source.Update()
	gauges := source.Gauges()
	counters := source.Counters()
	assert.Equal(t, len(counters), 1)
	assert.Equal(t, len(gauges), 28)
}

func TestPSUtilSource(t *testing.T) {
	source := NewPSUtilSource()
	source.Update()
	gauges := source.Gauges()
	counters := source.Counters()
	//! Число метрик зависит от числа процессоров но как минимум одна должна быть
	assert.Equal(t, len(gauges) > 0, true)
	assert.Equal(t, len(counters), 0)
}
