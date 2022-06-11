package storage

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMapRepository Тестируем вставку и чтение в MapRepository для gauge
func TestMapRepositoryGauge(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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

// TestMarshalJSON сериализация в json
func TestMarshalJSON(t *testing.T) {
	t.Parallel()

	repo := NewRepository()

	err := repo.UpdateGauge("testGaugeMetric", 100)
	assert.NoError(t, err)
	err = repo.UpdateCounter("testCounterMetric", 120)
	assert.NoError(t, err)

	jsRepo, err := json.Marshal(repo)
	assert.NoError(t, err)

	assert.Equal(t, string(jsRepo), `[{"id":"testGaugeMetric","type":"gauge","value":100},{"id":"testCounterMetric","type":"counter","delta":120}]`)
}

// TestUnmarshalJSON тест десериализации из json
func TestUnmarshalJSON(t *testing.T) {
	t.Parallel()

	repo := NewRepository()
	data := []byte(`[{"id":"testGaugeMetric","type":"gauge","value":100},{"id":"testCounterMetric","type":"counter","delta":120}]`)

	err := json.Unmarshal(data, &repo)
	assert.NoError(t, err)

	val, ok := repo.Gauge("testGaugeMetric")
	assert.Equal(t, ok, true)
	assert.Equal(t, val, float64(100))

	delta, ok := repo.Counter("testCounterMetric")
	assert.Equal(t, ok, true)
	assert.Equal(t, delta, int64(120))
}
