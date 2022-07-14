package config

import (
	"time"

	"github.com/ncyellow/devops/internal/gcfg"
)

// Config конфигурационные параметры сервера.
type Config struct {
	gcfg.GeneralConfig
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
	DatabaseConn  string        `env:"DATABASE_DSN"`
}
