// Package config определяет параметры настройки агента
package config

import (
	"time"

	"github.com/ncyellow/devops/internal/genconfig"
)

// Config содержит параметры по настройке агента
type Config struct {
	genconfig.GeneralConfig
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
}
