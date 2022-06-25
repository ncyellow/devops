package config

import "time"

// Config содержит параметры по настройке агента
type Config struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	SecretKey      string        `env:"KEY"`
}
