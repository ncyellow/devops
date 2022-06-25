package main

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/ncyellow/devops/internal/server"
	"github.com/ncyellow/devops/internal/server/config"
)

func main() {

	var cfg config.Config

	flag.StringVar(&cfg.Address, "a", "localhost:8080", "address in the format host:port")
	flag.DurationVar(&cfg.StoreInterval, "i", time.Second*300, "store interval in the format 300s")
	flag.BoolVar(&cfg.Restore, "r", true, "restore from file. true if needed")
	flag.StringVar(&cfg.StoreFile, "f", "/tmp/devops-metrics-db.json", "filename that used for save metrics state")
	flag.StringVar(&cfg.SecretKey, "k", "127.0.0.1:8080", "key for hash metrics")

	// Сначала парсим командную строку
	flag.Parse()

	// Далее приоритетно аргументы из ENV
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	server := server.Server{Conf: &cfg}
	server.RunServer()
}
