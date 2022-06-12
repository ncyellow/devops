package config

import "time"

type Config struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"` // "d:\\src\\golang\\devops-metrics-db.json"
	Restore       bool          `env:"RESTORE"`
}
