package main

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/ncyellow/devops/internal/agent"
)

func main() {

	var cfg agent.Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	collector := agent.Agent{Conf: cfg}
	log.Fatal(collector.Run())
}
