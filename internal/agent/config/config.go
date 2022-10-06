// Package config определяет параметры настройки агента
package config

import (
	"github.com/ncyellow/devops/internal/genconfig"
)

// Config содержит параметры по настройке агента
type Config struct {
	genconfig.GeneralConfig
	ReportInterval genconfig.Duration `env:"REPORT_INTERVAL" json:"report_interval"`
	PollInterval   genconfig.Duration `env:"POLL_INTERVAL" json:"poll_interval"`
}
