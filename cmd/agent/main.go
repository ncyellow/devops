package main

import (
	"github.com/ncyellow/devops/internal/agent"
	"log"
)

func main() {
	conf := agent.Config{
		Host: "localhost:8080",
	}
	collector := agent.Agent{Conf: conf}
	log.Fatal(collector.Run())
}
