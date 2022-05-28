package main

import (
	"log"

	"github.com/ncyellow/devops/internal/agent"
)

func main() {
	conf := agent.Config{
		Host: "localhost:8080",
	}
	collector := agent.Agent{Conf: conf}
	log.Fatal(collector.Run())
}
