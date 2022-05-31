package storage

import "fmt"

const (
	Gauge   = "gauge"
	Counter = "counter"
)

// Repository содержит API для работы с хранилищем метрик.
type Repository interface {
	// UpdateGauge обновить значение метрики типа gauge
	UpdateGauge(name string, value float64) error
	// UpdateCounter обновить значение метрики типа counter
	UpdateCounter(name string, value int64) error

	// Gauge возвращает текущее значение метрики типа gauge
	Gauge(name string) (val float64, ok bool)
	// Counter возвращает текущее значение метрики типа counter
	Counter(name string) (val int64, ok bool)

	// Stringer Вывод в строку всех метрик хранилища
	fmt.Stringer
}
