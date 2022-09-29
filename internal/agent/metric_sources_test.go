package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRuntimeSource проверяем что источник возвращает корректное число метрик
// counter и gauge
func TestRuntimeSource(t *testing.T) {
	source := RuntimeSource{}
	source.Update()
	gauges := source.Gauges()
	counters := source.Counters()
	//! Проверяем что приращение счетчика работает корректно
	assert.Equal(t, counters["PollCount"], int64(1))
	//! Проверяем что метрики gauges присутствуют,
	// но договорились что проверять наличие всех метрик избыточно
	assert.Equal(t, len(gauges) > 0, true)

	source.Update()
	gauges = source.Gauges()
	counters = source.Counters()
	//! После обновления метрик счетчик увеличился
	assert.Equal(t, counters["PollCount"], int64(2))
	//! Проверяем что метрики gauges присутствуют,
	// но договорились что проверять наличие всех метрик избыточно
	assert.Equal(t, len(gauges) > 0, true)
}

func TestPSUtilSource(t *testing.T) {
	source := NewPSUtilSource()
	source.Update()
	gauges := source.Gauges()
	counters := source.Counters()
	//! Число метрик зависит от числа процессоров, но как минимум одна должна быть
	assert.Equal(t, len(gauges) > 0, true)
	//! Метрики типа counters отсутствуют у источника данных, PSUtils
	assert.Equal(t, len(counters), 0)
}
