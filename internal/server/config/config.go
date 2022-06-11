package config

import "time"

type Config struct {
	Address       string        `env:"ADDRESS" envDefault:"localhost:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"` // "d:\\src\\golang\\devops-metrics-db.json"
	Restore       bool          `env:"RESTORE" envDefault:"true"`
}
