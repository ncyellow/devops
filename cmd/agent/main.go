package main

import (
	"github.com/ncyellow/devops/internal/agent"
	"log"
)

func main() {
	collector := agent.Agent{}
	log.Fatal(collector.Run())
}
