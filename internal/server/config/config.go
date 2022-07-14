package config

import (
	"time"

	"github.com/ncyellow/devops/internal/config"
)

// Config конфигурационные параметры сервера.
type Config struct {
	config.GeneralConfig
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
	DatabaseConn  string        `env:"DATABASE_DSN"`
}
