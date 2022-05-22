package storage

type Repository interface {
	UpdateGauge(name string, value float64) error
	UpdateCounter(name string, value int64) error

	Gauge(name string) (float64, bool)
	Counter(name string) (int64, bool)
}
