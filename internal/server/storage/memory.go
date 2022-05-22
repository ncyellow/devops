package storage

import (
	"fmt"
	"sort"
)

type MapRepository struct {
	gauges           map[string]float64
	counters         map[string]int64
	availableMetrics []string
}

func NewRepository() Repository {
	repo := MapRepository{}
	repo.gauges = make(map[string]float64)
	repo.counters = make(map[string]int64)

	//! Вспомогательная структура со списком допустимых метрик, если идет попытка добавить метрику не из этого
	//! списка - это ошибка
	repo.availableMetrics = []string{
		"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse",
		"HeapObjects", "HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse",
		"MSpanSys", "Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse",
		"StackSys", "Sys", "TotalAlloc", "RandomValue",
	}
	sort.Strings(repo.availableMetrics)
	return &repo
}

func (s *MapRepository) checkAvailable(name string) bool {
	i := sort.Search(len(s.availableMetrics),
		func(i int) bool { return s.availableMetrics[i] >= name })

	return i < len(s.availableMetrics) && s.availableMetrics[i] == name
}

func (s *MapRepository) UpdateGauge(name string, value float64) error {

	if s.checkAvailable(name) {
		s.gauges[name] = value
		return nil
	}
	return fmt.Errorf("metric with name %s is not available", name)
}

func (s *MapRepository) UpdateCounter(name string, value int64) error {
	if s.checkAvailable(name) {
		s.counters[name] = s.counters[name] + value
		return nil
	}
	return fmt.Errorf("metric with name %s is not available", name)
}

func (s *MapRepository) Gauge(name string) (val float64, ok bool) {
	val, ok = s.gauges[name]
	return
}

func (s *MapRepository) Counter(name string) (val int64, ok bool) {
	val, ok = s.counters[name]
	return
}
