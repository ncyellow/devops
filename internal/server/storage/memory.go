package storage

import "fmt"

// MapRepository хранилище метрик на основе map, реализует интерфейс Repository
type MapRepository struct {
	gauges   map[string]float64
	counters map[string]int64
}

// NewRepository конструктор
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

	gaugesText := ""
	for name, value := range s.gauges {
		gaugesText += fmt.Sprintf("<li>%s : %.3f</li>\n", name, value)
	}
	countersText := ""
	for name, value := range s.counters {
		countersText += fmt.Sprintf("<li>%s : %d</li>\n", name, value)
	}

	return fmt.Sprintf(htmlTmpl, gaugesText, countersText)
}
