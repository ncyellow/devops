package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
