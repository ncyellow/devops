package main

import (
	"flag"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/caarlos0/env/v6"
	"github.com/ncyellow/devops/internal/server"
	"github.com/ncyellow/devops/internal/server/config"
)

func main() {
	log.Info().Msg("Старт сервера")

	var cfg config.Config

	flag.StringVar(&cfg.Address, "a", "localhost:8080", "address in the format host:port")
	flag.DurationVar(&cfg.StoreInterval, "i", time.Second*300, "store interval in the format 300s")
	flag.BoolVar(&cfg.Restore, "r", true, "restore from file. true if needed")
	flag.StringVar(&cfg.StoreFile, "f", "/tmp/devops-metrics-db.json", "filename that used for save metrics state")
	flag.StringVar(&cfg.SecretKey, "k", "", "key for hash metrics")
	flag.StringVar(&cfg.DatabaseConn, "d", "", "connection string to postgresql")

	//flag.DurationVar(&cfg.StoreInterval, "i", time.Second*3, "store interval in the format 300s")
	//flag.StringVar(&cfg.DatabaseConn, "d", "user=postgres password=12345 host=localhost port=5433 dbname=gopractice", "connection string to postgresql")

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
