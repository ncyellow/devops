package main

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/ncyellow/devops/internal/server"
	"github.com/ncyellow/devops/internal/server/config"
)

func main() {

	var cfg config.Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	server := server.Server{Conf: cfg}
	server.RunServer()
}
