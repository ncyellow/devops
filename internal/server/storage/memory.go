package storage

type MapRepository struct {
	gauges           map[string]float64
	counters         map[string]int64
	availableMetrics []string
}

func NewRepository() Repository {
	repo := MapRepository{}
	repo.gauges = make(map[string]float64)
	repo.counters = make(map[string]int64)
	return &repo
}

func (s *MapRepository) UpdateGauge(name string, value float64) error {

	s.gauges[name] = value
	return nil
}

func (s *MapRepository) UpdateCounter(name string, value int64) error {
	s.counters[name] = s.counters[name] + value
	return nil
}

func (s *MapRepository) Gauge(name string) (val float64, ok bool) {
	val, ok = s.gauges[name]
	return
}

func (s *MapRepository) Counter(name string) (val int64, ok bool) {
	val, ok = s.counters[name]
	return
}
