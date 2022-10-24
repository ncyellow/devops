package repository

import (
	"encoding/json"

	"testing"

	"github.com/ncyellow/devops/internal/genconfig"
	"github.com/stretchr/testify/assert"
)

// TestMapRepository Тестируем вставку и чтение в MapRepository для gauge
func TestMapRepositoryGauge(t *testing.T) {
	t.Parallel()

	repo := NewRepository(&genconfig.GeneralConfig{})

	// обновление
	repo.UpdateGauge("testGauge", 100.0)

	// чтение
	val, ok := repo.Gauge("testGauge")
	assert.Equal(t, 100.0, val)
	assert.Equal(t, true, ok)

	// обновляем повторно
	repo.UpdateGauge("testGauge", 300.0)

	// проверяем что старое значение перезаписалось
	val, ok = repo.Gauge("testGauge")
	assert.Equal(t, 300.0, val)
	assert.Equal(t, true, ok)

	// Проверка чтения неизвестного значения
	_, ok = repo.Gauge("unknownGauge")
	assert.Equal(t, false, ok)
}

// TestMapRepository Тестируем вставку и чтение в MapRepository для counter
func TestMapRepositoryCounter(t *testing.T) {
	t.Parallel()

	repo := NewRepository(&genconfig.GeneralConfig{})

	// обновление
	repo.UpdateCounter("testCounter", 100)

	// чтение
	val, ok := repo.Counter("testCounter")
	assert.Equal(t, int64(100), val)
	assert.Equal(t, true, ok)

	// обновляем еще раз
	repo.UpdateCounter("testCounter", 100)

	// проверяем что счетчик приплюсовал значение
	val, ok = repo.Counter("testCounter")
	assert.Equal(t, int64(200), val)
	assert.Equal(t, true, ok)

	// Проверка чтения неизвестного значения
	_, ok = repo.Counter("unknownCounter")
	assert.Equal(t, false, ok)
}

// TestMapRepository Тестируем вставку и чтение в MapRepository для counter
func TestMapRepositoryMetricsCounter(t *testing.T) {
	t.Parallel()

	repo := NewRepository(&genconfig.GeneralConfig{})

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
	assert.Equal(t, updateValue, *val.Delta)
	assert.Equal(t, true, ok)

	// обновляем еще раз
	err = repo.UpdateMetric(Metrics{
		ID:    "testCounterMetric",
		MType: Counter,
		Delta: &updateValue,
	})
	assert.NoError(t, err)

	// чтение
	val, ok = repo.Metric("testCounterMetric", Counter)
	assert.Equal(t, updateValue*2, *val.Delta)
	assert.Equal(t, true, ok)

	// Проверка чтения неизвестного значения
	_, ok = repo.Metric("unknownMetricCoutner", Counter)
	assert.Equal(t, ok, false)
}

// TestMapRepository Тестируем вставку и чтение в MapRepository для counter
func TestMapRepositoryStringer(t *testing.T) {
	t.Parallel()

	repo := NewRepository(&genconfig.GeneralConfig{})

	// обновление
	repo.UpdateGauge("testGauge", 100.0)

	repo.UpdateCounter("testCounter", 100)

	correctHTML := `
	<html>
	<body>
	<h1>All metrics</h1>
	<h3>gauges</h3>
	<ul>
	  <li>testGauge : 100.000</li>

	</ul>
	<h3>counters</h3>
	<ul>
	  <li>testCounter : 100</li>

	</ul>
	</body>
	</html>`

	assert.Equal(t, correctHTML, RenderHTML(repo.ToMetrics()))

}

// TestMapRepository Тестируем вставку и чтение в MapRepository для gauge
func TestMapRepositoryMetricsGauge(t *testing.T) {
	t.Parallel()

	repo := NewRepository(&genconfig.GeneralConfig{})

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
	assert.Equal(t, updateValue, *val.Value)
	assert.Equal(t, true, ok)

	// обновляем повторно
	updateValue = 300
	err = repo.UpdateMetric(Metrics{
		ID:    "testGaugeMetric",
		MType: Gauge,
		Value: &updateValue,
	})
	assert.NoError(t, err)

	// обновляем с левым типом
	updateValue = 300
	err = repo.UpdateMetric(Metrics{
		ID:    "testMetric",
		MType: "unknownType",
		Value: &updateValue,
	})
	assert.Error(t, err)

	// проверяем что старое значение перезаписалось
	val, ok = repo.Metric("testGaugeMetric", Gauge)
	assert.Equal(t, updateValue, *val.Value)
	assert.Equal(t, true, ok)

	// Проверка чтения неизвестной метрики тика Gauge
	val, ok = repo.Metric("unknownMetricGauge", Gauge)
	assert.Equal(t, false, ok)
	assert.Equal(t, val, Metrics{})

	// Проверка чтения неизвестной метрики тика Counter
	val, ok = repo.Metric("unknownMetricCounter", Counter)
	assert.Equal(t, false, ok)
	assert.Equal(t, val, Metrics{})

	// Проверка чтения неизвестной метрики тика Counter
	val, ok = repo.Metric("unknownMetric", "unknownType")
	assert.Equal(t, false, ok)
	assert.Equal(t, val, Metrics{})
}

// TestMarshalJSON сериализация в json
func TestMarshalJSON(t *testing.T) {
	t.Parallel()

	repo := NewRepository(&genconfig.GeneralConfig{})

	repo.UpdateGauge("testGaugeMetric", 100)

	repo.UpdateCounter("testCounterMetric", 120)

	jsRepo, err := json.Marshal(repo)
	assert.NoError(t, err)

	assert.JSONEq(t, string(jsRepo), `[{"id":"testGaugeMetric","type":"gauge","value":100},{"id":"testCounterMetric","type":"counter","delta":120}]`)
}

// TestUnmarshalJSON тест десериализации из json
func TestUnmarshalJSON(t *testing.T) {
	t.Parallel()

	repo := NewRepository(&genconfig.GeneralConfig{})

	data := []byte(`[{"id":"testGaugeMetric","type":"gauge","value":100},{"id":"testCounterMetric","type":"counter","delta":120}]`)

	err := json.Unmarshal(data, &repo)
	assert.NoError(t, err)

	val, ok := repo.Gauge("testGaugeMetric")
	assert.Equal(t, true, ok)
	assert.Equal(t, float64(100), val)

	delta, ok := repo.Counter("testCounterMetric")
	assert.Equal(t, true, ok)
	assert.Equal(t, int64(120), delta)

	brokenRepo := NewRepository(&genconfig.GeneralConfig{})
	brokenData := []byte(`{"name": "Joe", "age": null, }`)

	err = json.Unmarshal(brokenData, &brokenRepo)
	assert.Error(t, err)
}
