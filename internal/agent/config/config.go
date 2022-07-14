package config

import (
	"time"

	"github.com/ncyellow/devops/internal/config"
)

// Config содержит параметры по настройке агента
type Config struct {
	config.GeneralConfig
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
}
