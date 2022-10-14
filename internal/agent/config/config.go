// Package config определяет параметры настройки агента
package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ncyellow/devops/internal/genconfig"
	"github.com/rs/zerolog/log"
)

// Config содержит параметры по настройке агента
type Config struct {
	genconfig.GeneralConfig
	ReportInterval genconfig.Duration `env:"REPORT_INTERVAL" json:"report_interval"`
	PollInterval   genconfig.Duration `env:"POLL_INTERVAL" json:"poll_interval"`
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
