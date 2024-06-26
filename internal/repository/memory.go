// Package repository содержит имплементацию Repositiory для хранения метрик в памяти
package repository

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/ncyellow/devops/internal/genconfig"
	"github.com/ncyellow/devops/internal/hash"
)

// MapRepository структура данных для работы метриками на основе map, реализует интерфейс Repository
type MapRepository struct {
	conf     *genconfig.GeneralConfig
	gauges   map[string]float64
	counters map[string]int64

	// Мьютекcа два так как обе структуры данных у нас независимы и так мы уменьшаем гранулярность блокировок
	gaugesLock   *sync.RWMutex
	countersLock *sync.RWMutex
}

// NewRepository конструктор
func NewRepository(conf *genconfig.GeneralConfig) Repository {
	repo := MapRepository{}
	repo.conf = conf
	repo.gauges = make(map[string]float64)
	repo.counters = make(map[string]int64)

	repo.gaugesLock = &sync.RWMutex{}
	repo.countersLock = &sync.RWMutex{}
	return &repo
}

func (s *MapRepository) UpdateGauge(name string, value float64) {
	s.gaugesLock.Lock()
	s.gauges[name] = value
	s.gaugesLock.Unlock()
}

func (s *MapRepository) UpdateCounter(name string, value int64) {
	s.countersLock.Lock()
	s.counters[name] = s.counters[name] + value
	s.countersLock.Unlock()
}

func (s *MapRepository) UpdateMetric(metric Metrics) error {
	switch metric.MType {
	case Gauge:
		s.UpdateGauge(metric.ID, *metric.Value)
	case Counter:
		s.UpdateCounter(metric.ID, *metric.Delta)
	default:
		return fmt.Errorf("metric with type %s doesn't exsist", metric.MType)
	}
	return nil
}

func (s *MapRepository) Gauge(name string) (val float64, ok bool) {
	s.gaugesLock.RLock()
	val, ok = s.gauges[name]
	s.gaugesLock.RUnlock()

	return
}

func (s *MapRepository) Counter(name string) (val int64, ok bool) {
	s.countersLock.RLock()
	val, ok = s.counters[name]
	s.countersLock.RUnlock()
	return
}

func (s *MapRepository) Metric(name string, mType string) (val Metrics, ok bool) {
	encodeFunc := hash.CreateEncodeFunc(s.conf.SecretKey)
	switch mType {
	case Gauge:
		val, ok := s.Gauge(name)
		if !ok {
			return Metrics{}, ok
		}
		metric := Metrics{
			ID:    name,
			MType: mType,
			Value: &val,
			Delta: nil,
		}
		metric.CalcHash(encodeFunc)
		return metric, ok
	case Counter:
		val, ok := s.Counter(name)
		if !ok {
			return Metrics{}, ok
		}
		metric := Metrics{
			ID:    name,
			MType: mType,
			Value: nil,
			Delta: &val,
		}
		metric.CalcHash(encodeFunc)
		return metric, ok
	default:
		return Metrics{}, false
	}
}

// ToMetrics Конвертация данных MapRepository в []Metrics
func (s *MapRepository) ToMetrics() []Metrics {
	totalCount := len(s.gauges) + len(s.counters)
	metrics := make([]Metrics, 0, totalCount)
	hashFunc := hash.CreateEncodeFunc(s.conf.SecretKey)

	s.gaugesLock.RLock()
	for name, value := range s.gauges {
		gaugeValue := value
		metric := Metrics{
			ID:    name,
			MType: Gauge,
			Value: &gaugeValue,
		}
		metric.Hash = metric.CalcHash(hashFunc)
		metrics = append(metrics, metric)
	}
	s.gaugesLock.RUnlock()

	s.countersLock.RLock()
	for name, value := range s.counters {
		counterValue := value
		metric := Metrics{
			ID:    name,
			MType: Counter,
			Delta: &counterValue,
		}
		metric.Hash = metric.CalcHash(hashFunc)
		metrics = append(metrics, metric)
	}
	s.countersLock.RUnlock()
	return metrics
}

// FromMetrics - обновляет метрики в MapRepository по []Metrics
func (s *MapRepository) FromMetrics(metrics []Metrics) {
	for _, metric := range metrics {
		switch metric.MType {
		case Gauge:
			if metric.Value != nil {
				s.UpdateGauge(metric.ID, *metric.Value)
			}
		case Counter:
			if metric.Delta != nil {
				s.UpdateCounter(metric.ID, *metric.Delta)
			}
		}
	}
}

// Clear - очищаем все метрики репозитория
func (s *MapRepository) Clear() {
	s.gaugesLock.Lock()
	s.gauges = make(map[string]float64)
	s.gaugesLock.Unlock()

	s.countersLock.Lock()
	s.counters = make(map[string]int64)
	s.countersLock.Unlock()
}

// MarshalJSON - реализация интерфейса Marshaler
func (s *MapRepository) MarshalJSON() ([]byte, error) {
	metrics := s.ToMetrics()
	jsMetrics, err := json.Marshal(metrics)
	if err != nil {
		return []byte{}, nil
	}
	return jsMetrics, nil
}

// UnmarshalJSON - реализация интерфейса Unmarshaler
func (s *MapRepository) UnmarshalJSON(data []byte) error {
	var metrics []Metrics
	err := json.Unmarshal(data, &metrics)
	if err != nil {
		return err
	}
	s.FromMetrics(metrics)
	return nil
}
