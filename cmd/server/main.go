package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/caarlos0/env/v6"
	"github.com/ncyellow/devops/internal/server"
	"github.com/ncyellow/devops/internal/server/config"
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

	log.Info().Msg("Старт сервера")

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

	flag.StringVar(&cfg.Address, "a", "localhost:8080", "address in the format host:port")
	flag.DurationVar(&cfg.StoreInterval.Duration, "i", time.Second*300, "store interval in the format 300s")
	flag.BoolVar(&cfg.Restore, "r", true, "restore from file. true if needed")
	flag.StringVar(&cfg.StoreFile, "f", "/tmp/devops-metrics-db.json", "filename that used for save metrics state")
	flag.StringVar(&cfg.SecretKey, "k", "", "key for hash metrics")
	flag.StringVar(&cfg.CryptoKey, "crypto-key", "", "private server crypto key")
	flag.StringVar(&cfg.DatabaseConn, "d", "", "connection string to postgresql")
	flag.StringVar(&cfg.TrustedSubNet, "t", "", "trusted subnet cidr")

	// Сначала парсим командную строку
	flag.Parse()

	// Далее приоритетно аргументы из ENV
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal().Err(err)
	}

	log.Info().Msgf("Настройки запуска сервера - %#v\n", cfg)

	server := server.Server{Conf: &cfg}
	server.RunServer()
}
