package agent

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// MetricSource интерфейс источника для сбора метрик
type MetricSource interface {
	// Update - читает метрики если есть необходимость выполняет агрегацию и объединение
	Update()
	// Counters - возвращает список метрик типа Counter
	Counters() map[string]int64
	// Gauges - возвращает список метрик типа Gauge
	Gauges() map[string]float64
}

// RuntimeSource реализация источника метрик на основании пакета runtime
type RuntimeSource struct {
	pollCount   int64
	randomValue float64
	runtime.MemStats
}

// PSUtilSource реализация источника метрик на основании пакета gopsutil
type PSUtilSource struct {
	gauges map[string]float64
}

func (rs *RuntimeSource) Update() {
	//! Обновляем все стандартные метрики
	//! Инкремент счетчика и новый рандом
	runtime.ReadMemStats(&rs.MemStats)
	rs.pollCount += 1

	rand.Seed(time.Now().UnixNano())
	rs.randomValue = rand.Float64()
}

func (rs *RuntimeSource) Counters() map[string]int64 {
	return map[string]int64{
		"PollCount": rs.pollCount,
	}
}

func (rs *RuntimeSource) Gauges() map[string]float64 {
	return map[string]float64{
		"Alloc":         float64(rs.Alloc),
		"BuckHashSys":   float64(rs.BuckHashSys),
		"Frees":         float64(rs.Frees),
		"GCCPUFraction": rs.GCCPUFraction,
		"GCSys":         float64(rs.GCSys),
		"HeapAlloc":     float64(rs.HeapAlloc),
		"HeapIdle":      float64(rs.HeapIdle),
		"HeapInuse":     float64(rs.HeapInuse),
		"HeapObjects":   float64(rs.HeapObjects),
		"HeapReleased":  float64(rs.HeapReleased),
		"HeapSys":       float64(rs.HeapSys),
		"LastGC":        float64(rs.LastGC),
		"Lookups":       float64(rs.Lookups),
		"MCacheInuse":   float64(rs.MCacheInuse),
		"MCacheSys":     float64(rs.MCacheSys),
		"MSpanInuse":    float64(rs.MSpanInuse),
		"MSpanSys":      float64(rs.MSpanSys),
		"Mallocs":       float64(rs.Mallocs),
		"NextGC":        float64(rs.NextGC),
		"NumForcedGC":   float64(rs.NumForcedGC),
		"NumGC":         float64(rs.NumGC),
		"OtherSys":      float64(rs.OtherSys),
		"PauseTotalNs":  float64(rs.PauseTotalNs),
		"StackInuse":    float64(rs.StackInuse),
		"StackSys":      float64(rs.StackSys),
		"Sys":           float64(rs.Sys),
		"TotalAlloc":    float64(rs.TotalAlloc),
		"RandomValue":   rs.randomValue,
	}
}

// NewPSUtilSource инициализация объекта PSUtilSource
func NewPSUtilSource() *PSUtilSource {
	source := PSUtilSource{}
	source.gauges = make(map[string]float64)
	return &source
}

func (ps *PSUtilSource) Update() {

	// Читаем метрики RAM
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Info().Msg("get memory metrics failed")
		return
	}

	ps.gauges["FreeMemory"] = float64(v.Free)
	ps.gauges["TotalMemory"] = float64(v.Total)

	// Читаем метрики cpu utilization
	percent, err := cpu.Percent(time.Second, true)
	if err != nil {
		log.Info().Msg("get cpu metrics failed")
	}
	for index, value := range percent {
		cpuUtil := value
		ps.gauges[fmt.Sprintf("CPUutilization%d", index)] = cpuUtil
	}

}

func (ps *PSUtilSource) Counters() map[string]int64 {
	return map[string]int64{}
}

func (ps *PSUtilSource) Gauges() map[string]float64 {
	return ps.gauges
}
