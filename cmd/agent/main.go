package main

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/ncyellow/devops/internal/agent"
)

func main() {

	var cfg agent.Config

	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "address in the format host:port")
	flag.DurationVar(&cfg.ReportInterval, "r", time.Second*10, "send to server interval in the format 10s")
	flag.DurationVar(&cfg.PollInterval, "p", time.Second*2, "polling metrics interval in the format 2s")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	collector := agent.Agent{Conf: cfg}
	log.Fatal(collector.Run())
}
