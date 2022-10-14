// Package config содержит настройки сервера
package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ncyellow/devops/internal/genconfig"
	"github.com/rs/zerolog/log"
)

// Config конфигурационные параметры сервера.
type Config struct {
	genconfig.GeneralConfig
	StoreInterval genconfig.Duration `env:"STORE_INTERVAL" json:"store_interval"`
	StoreFile     string             `env:"STORE_FILE" json:"store_file"`
	Restore       bool               `env:"RESTORE" json:"restore"`
	DatabaseConn  string             `env:"DATABASE_DSN" json:"database_dsn"`
	TrustedSubNet string             `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
}

func ReadConfig(fileName string) Config {
	cfg := Config{}
	if fileName != "" {
		file, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Info().Msgf("Ошибка при чтении конфигурационного файла %s", err.Error())
			return cfg
		}
		err = json.Unmarshal(file, &cfg)
		if err != nil {
			log.Info().Msgf("Ошибка при разборе конфигурационного файла %s", err.Error())
		}
	}
	return cfg
}
