package config

import "time"

// Config конфигурационные параметры сервера.
type Config struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
	SecretKey     string        `env:"KEY"`
	DatabaseConn  string        `env:"DATABASE_DSN"`
}
