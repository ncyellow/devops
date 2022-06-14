package storage

import (
	"encoding/json"
	"fmt"
	"sync"
)

// MapRepository хранилище метрик на основе map, реализует интерфейс Repository
type MapRepository struct {
	gauges   map[string]float64
	counters map[string]int64

	// Мьютекcа два так как обе структуры данных у нас независимы и так мы уменьшаем гранулярность блокировок
	gaugesLock   *sync.RWMutex
	countersLock *sync.RWMutex
}

// NewRepository конструктор
func NewRepository() Repository {
	repo := MapRepository{}
	repo.gauges = make(map[string]float64)
	repo.counters = make(map[string]int64)

	repo.gaugesLock = &sync.RWMutex{}
	repo.countersLock = &sync.RWMutex{}
	return &repo
}

func (s *MapRepository) UpdateGauge(name string, value float64) error {
	s.gaugesLock.Lock()
	s.gauges[name] = value
	s.gaugesLock.Unlock()
	return nil
}

func (s *MapRepository) UpdateCounter(name string, value int64) error {
	s.countersLock.Lock()
	s.counters[name] = s.counters[name] + value
	s.countersLock.Unlock()
	return nil
}

func (s *MapRepository) UpdateMetric(metric Metrics) error {
	if metric.MType == Gauge {
		return s.UpdateGauge(metric.ID, *metric.Value)
	} else if metric.MType == Counter {
		return s.UpdateCounter(metric.ID, *metric.Delta)
	}
	return fmt.Errorf("metric with type %s doesn't exsist", metric.MType)
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
	if mType == Gauge {
		val, ok := s.Gauge(name)
		if !ok {
			return Metrics{}, ok
		}
		return Metrics{
			ID:    name,
			MType: mType,
			Value: &val,
			Delta: nil,
		}, ok
	}
	if mType == Counter {
		val, ok := s.Counter(name)
		if !ok {
			return Metrics{}, ok
		}
		return Metrics{
			ID:    name,
			MType: mType,
			Value: nil,
			Delta: &val,
		}, ok
	}
	return Metrics{}, false
}

func (s *MapRepository) String() string {
	htmlTmpl := `
	<html>
	<body>
	<h1>All metrics</h1>
	<h3>gauges</h3>
	<ul>
	  %s
	</ul>
	<h3>counters</h3>
	<ul>
	  %s
	</ul>
	</body>
	</html>`

	s.gaugesLock.RLock()
	gaugesText := ""
	for name, value := range s.gauges {
		gaugesText += fmt.Sprintf("<li>%s : %.3f</li>\n", name, value)
	}
	s.gaugesLock.RUnlock()

	s.countersLock.RLock()
	countersText := ""
	for name, value := range s.counters {
		countersText += fmt.Sprintf("<li>%s : %d</li>\n", name, value)
	}
	s.countersLock.RUnlock()

	return fmt.Sprintf(htmlTmpl, gaugesText, countersText)
}

// toMetrics Конвертация данных MapRepository в []Metrics
func (s *MapRepository) toMetrics() []Metrics {
	totalCount := len(s.gauges) + len(s.counters)
	metrics := make([]Metrics, 0, totalCount)

	s.gaugesLock.RLock()
	for name, value := range s.gauges {
		gaugeValue := value
		metrics = append(metrics, Metrics{
			ID:    name,
			MType: Gauge,
			Value: &gaugeValue,
		})
	}
	s.gaugesLock.RUnlock()

	s.countersLock.RLock()
	for name, value := range s.counters {
		counterValue := value
		metrics = append(metrics, Metrics{
			ID:    name,
			MType: Counter,
			Delta: &counterValue,
		})
	}
	s.countersLock.RUnlock()
	return metrics
}

// fromMetrics - обновляет метрики в MapRepository по []Metrics
func (s *MapRepository) fromMetrics(metrics []Metrics) {
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

// MarshalJSON - реализация интерфейса Marshaler
func (s *MapRepository) MarshalJSON() ([]byte, error) {
	metrics := s.toMetrics()
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
	s.fromMetrics(metrics)
	return nil
}
