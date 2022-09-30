// Package repository содержит интерфейс для хранения метрик и сам тип метрики
package repository

import (
	"fmt"

	"github.com/ncyellow/devops/internal/hash"
)

const (
	Gauge   = "gauge"
	Counter = "counter"
)

// Metrics тип метрики для взаимодействия по сети и хранения на диске
type Metrics struct {
	// Имя метрики
	ID string `json:"id"`
	// Параметр, принимающий значение gauge или counter
	MType string `json:"type"`
	// Значение метрики в случае передачи counter
	Delta *int64 `json:"delta,omitempty"`
	// Значение метрики в случае передачи gauge
	Value *float64 `json:"value,omitempty"`
	// Значение хеш-функции
	Hash string `json:"hash,omitempty"`
}

// CalcHash вычисление хеша с подписью метрики
func (m *Metrics) CalcHash(encodeFunc hash.EncodeFunc) string {
	switch m.MType {
	case Gauge:
		return encodeFunc(fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value))
	case Counter:
		return encodeFunc(fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta))
	default:
		return ""
	}
}

// Repository содержит API для работы с метриками.
// Хранение разделено на две сущности. Кеш в RAM - Repository. А PersistentStorage
// представляет сохранение в долговременное хранилище файл или бд
type Repository interface {
	// UpdateGauge обновить значение метрики типа gauge
	UpdateGauge(name string, value float64) error
	// UpdateCounter обновить значение метрики типа counter
	UpdateCounter(name string, value int64) error

	// Gauge возвращает текущее значение метрики типа gauge
	Gauge(name string) (val float64, ok bool)
	// Counter возвращает текущее значение метрики типа counter
	Counter(name string) (val int64, ok bool)

	// Metric возвращает значение метрики по названию
	Metric(name string, mType string) (val Metrics, ok bool)

	// UpdateMetric обновляет данные в хранилище по значению Metrics
	UpdateMetric(metrics Metrics) error

	// FromMetrics загрузить данные в репозиторий из []Metrics
	FromMetrics(metrics []Metrics)

	// ToMetrics экспорт данных репозитория в []Metrics
	ToMetrics() []Metrics
}
