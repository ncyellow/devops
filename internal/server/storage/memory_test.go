package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMapRepository Тестируем вставку и чтение в MapRepository для gauge
func TestMapRepositoryGauge(t *testing.T) {
	repo := NewRepository()

	// обновление
	err := repo.UpdateGauge("testGauge", 100.0)
	assert.NoError(t, err)

	// чтение
	val, ok := repo.Gauge("testGauge")
	assert.Equal(t, val, 100.0)
	assert.Equal(t, ok, true)

	// обновляем повторно
	err = repo.UpdateGauge("testGauge", 300.0)
	assert.NoError(t, err)

	// проверяем что старое значение перезаписалось
	val, ok = repo.Gauge("testGauge")
	assert.Equal(t, val, 300.0)
	assert.Equal(t, ok, true)

	// Проверка чтения неизвестного значения
	_, ok = repo.Gauge("unknownGauge")
	assert.Equal(t, ok, false)
}

// TestMapRepository Тестируем вставку и чтение в MapRepository для counter
func TestMapRepositoryCounter(t *testing.T) {
	repo := NewRepository()

	// обновление
	err := repo.UpdateCounter("testCounter", 100)
	assert.NoError(t, err)

	// чтение
	val, ok := repo.Counter("testCounter")
	assert.Equal(t, val, int64(100))
	assert.Equal(t, ok, true)

	// обновляем еще раз
	err = repo.UpdateCounter("testCounter", 100)
	assert.NoError(t, err)

	// проверяем что счетчик приплюсовал значение
	val, ok = repo.Counter("testCounter")
	assert.Equal(t, val, int64(200))
	assert.Equal(t, ok, true)

	// Проверка чтения неизвестного значения
	_, ok = repo.Counter("unknownCounter")
	assert.Equal(t, ok, false)
}

// TestMapRepository Тестируем вставку и чтение в MapRepository для counter
func TestMapRepositoryMetricsCounter(t *testing.T) {
	repo := NewRepository()

	var updateValue int64 = 100

	// обновление
	err := repo.UpdateMetric(Metrics{
		ID:    "testCounterMetric",
		MType: Counter,
		Delta: &updateValue,
	})
	assert.NoError(t, err)

	// чтение
	val, ok := repo.Metric("testCounterMetric", Counter)
	assert.Equal(t, *val.Delta, updateValue)
	assert.Equal(t, ok, true)

	// обновляем еще раз
	err = repo.UpdateMetric(Metrics{
		ID:    "testCounterMetric",
		MType: Counter,
		Delta: &updateValue,
	})
	assert.NoError(t, err)

	// чтение
	val, ok = repo.Metric("testCounterMetric", Counter)
	assert.Equal(t, *val.Delta, updateValue*2)
	assert.Equal(t, ok, true)

	// Проверка чтения неизвестного значения
	_, ok = repo.Metric("unknownMetricCoutner", Counter)
	assert.Equal(t, ok, false)
}

// TestMapRepository Тестируем вставку и чтение в MapRepository для gauge
func TestMapRepositoryMetricsGauge(t *testing.T) {
	repo := NewRepository()

	// обновление
	var updateValue float64 = 100

	// обновление
	err := repo.UpdateMetric(Metrics{
		ID:    "testGaugeMetric",
		MType: Gauge,
		Value: &updateValue,
	})
	assert.NoError(t, err)

	// чтение
	val, ok := repo.Metric("testGaugeMetric", Gauge)
	assert.Equal(t, *val.Value, updateValue)
	assert.Equal(t, ok, true)

	// обновляем повторно
	updateValue = 300
	err = repo.UpdateMetric(Metrics{
		ID:    "testGaugeMetric",
		MType: Gauge,
		Value: &updateValue,
	})
	assert.NoError(t, err)

	// проверяем что старое значение перезаписалось
	val, ok = repo.Metric("testGaugeMetric", Gauge)
	assert.Equal(t, *val.Value, updateValue)
	assert.Equal(t, ok, true)

	// Проверка чтения неизвестного значения
	_, ok = repo.Metric("unknownMetricGauge", Counter)
	assert.Equal(t, ok, false)
}
