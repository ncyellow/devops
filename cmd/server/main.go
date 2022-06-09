package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v6"
	"github.com/ncyellow/devops/internal/server/handlers"
	"github.com/ncyellow/devops/internal/server/storage"
)

type Config struct {
	Address string `env:"ADDRESS" envDefault:":8080"`
}

func main() {

	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	repo := storage.NewRepository()
	r := handlers.NewRouter(repo)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := http.ListenAndServe(cfg.Address, r); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	<-done
}
