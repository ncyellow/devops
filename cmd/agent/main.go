package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/caarlos0/env/v6"
	"github.com/ncyellow/devops/internal/agent"
	"github.com/ncyellow/devops/internal/agent/config"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	log.Info().Msg("Старт агента")

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	var cfg config.Config
	confEnvFile := os.Getenv("CONFIG")
	// Сначала смотрим задан ли конфиг в env
	if confEnvFile != "" {
		cfg = config.ReadConfig(confEnvFile)
	} else {
		// Раз не задан пытаемся нащупать его в параметрах
		confFile := flag.String("c", "", "config file")
		flag.Parse()

		//! Начинаем с меньшего приоритета
		if *confFile != "" {
			fmt.Println(*confFile)
			cfg = config.ReadConfig(*confFile)
		}
	}

	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "address in the format host:port")
	flag.DurationVar(&cfg.ReportInterval.Duration, "r", time.Second*10, "send to server interval in the format 10s")
	flag.DurationVar(&cfg.PollInterval.Duration, "p", time.Second*2, "polling metrics interval in the format 2s")
	flag.StringVar(&cfg.SecretKey, "k", "", "key for hash metrics")
	flag.StringVar(&cfg.CryptoKey, "crypto-key", "", "public agent crypto key")

	// Сначала аргументы командной строки
	flag.Parse()

	// Далее более приоритетные от ENV
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal().Err(err)
	}
	log.Info().Msgf("Настройки запуска агента - %#v\n", cfg)

	collector := agent.Agent{Conf: &cfg}
	log.Fatal().Err(collector.Run())
}
