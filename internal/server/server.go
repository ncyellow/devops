package server

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ncyellow/devops/internal/server/config"
	"github.com/ncyellow/devops/internal/server/handlers"
	"github.com/ncyellow/devops/internal/server/storage"
)

type Server struct {
	Conf config.Config
}

func (s Server) RunServer() {
	repo := storage.NewRepository()

	if s.Conf.Restore {
		storage.RestoreFromFile(s.Conf.StoreFile, repo)
	}

	r := handlers.NewRouter(repo)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := http.ListenAndServe(s.Conf.Address, r); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	go storage.RunStorageSaver(s.Conf, repo)

	<-done
}
