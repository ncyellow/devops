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

type MetricSource interface {
	Update()
	Counters() map[string]int64
	Gauges() map[string]float64
}

// RuntimeMetrics текущее состояние всех метрик обновляются с интервалом pollInterval
type RuntimeMetrics struct {
	PollCount   int64
	RandomValue float64
	runtime.MemStats
}

type RuntimeSource struct {
	metrics RuntimeMetrics
}

func (rs *RuntimeSource) Update() {
	//! Обновляем все стандартные метрики
	//! Инкремент счетчика и новый рандом
	runtime.ReadMemStats(&rs.metrics.MemStats)
	rs.metrics.PollCount += 1

	rand.Seed(time.Now().UnixNano())
	rs.metrics.RandomValue = rand.Float64()
}

func (rs *RuntimeSource) Counters() map[string]int64 {
	return map[string]int64{
		"PollCount": rs.metrics.PollCount,
	}
}

func (rs *RuntimeSource) Gauges() map[string]float64 {
	return map[string]float64{
		"Alloc":         float64(rs.metrics.Alloc),
		"BuckHashSys":   float64(rs.metrics.BuckHashSys),
		"Frees":         float64(rs.metrics.Frees),
		"GCCPUFraction": rs.metrics.GCCPUFraction,
		"GCSys":         float64(rs.metrics.GCSys),
		"HeapAlloc":     float64(rs.metrics.HeapAlloc),
		"HeapIdle":      float64(rs.metrics.HeapIdle),
		"HeapInuse":     float64(rs.metrics.HeapInuse),
		"HeapObjects":   float64(rs.metrics.HeapObjects),
		"HeapReleased":  float64(rs.metrics.HeapReleased),
		"HeapSys":       float64(rs.metrics.HeapSys),
		"LastGC":        float64(rs.metrics.LastGC),
		"Lookups":       float64(rs.metrics.Lookups),
		"MCacheInuse":   float64(rs.metrics.MCacheInuse),
		"MCacheSys":     float64(rs.metrics.MCacheSys),
		"MSpanInuse":    float64(rs.metrics.MSpanInuse),
		"MSpanSys":      float64(rs.metrics.MSpanSys),
		"Mallocs":       float64(rs.metrics.Mallocs),
		"NextGC":        float64(rs.metrics.NextGC),
		"NumForcedGC":   float64(rs.metrics.NumForcedGC),
		"NumGC":         float64(rs.metrics.NumGC),
		"OtherSys":      float64(rs.metrics.OtherSys),
		"PauseTotalNs":  float64(rs.metrics.PauseTotalNs),
		"StackInuse":    float64(rs.metrics.StackInuse),
		"StackSys":      float64(rs.metrics.StackSys),
		"Sys":           float64(rs.metrics.Sys),
		"TotalAlloc":    float64(rs.metrics.TotalAlloc),
		"RandomValue":   rs.metrics.RandomValue,
	}
}

type PSUtilSource struct {
	gauges map[string]float64
}

func NewPSUtilSource() *PSUtilSource {
	source := PSUtilSource{}
	source.gauges = make(map[string]float64)
	return &source
}

func (ps *PSUtilSource) Update() {
	//! Обновляем все стандартные метрики
	//! Инкремент счетчика и новый рандом
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Info().Msg("get memory metrics failed")
		return
	}

	ps.gauges["FreeMemory"] = float64(v.Free)
	ps.gauges["TotalMemory"] = float64(v.Total)

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
