// Package config содержит настройки сервера
package config

import (
	"github.com/ncyellow/devops/internal/genconfig"
)

// Config конфигурационные параметры сервера.
type Config struct {
	genconfig.GeneralConfig
	StoreInterval genconfig.Duration `env:"STORE_INTERVAL" json:"store_interval"`
	StoreFile     string             `env:"STORE_FILE" json:"store_file"`
	Restore       bool               `env:"RESTORE" json:"restore"`
	DatabaseConn  string             `env:"DATABASE_DSN" json:"database_dsn"`
}
