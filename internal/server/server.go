package server

import (
	"fmt"
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
	Conf *config.Config
}

func (s Server) RunServer() {
	repo := storage.NewRepository(s.Conf)

	saver, err := storage.CreateSaver(s.Conf)
	if err != nil {
		fmt.Println("cant create NewSaver")
	}
	defer saver.Close(repo)

	saver.Load(repo)

	r := handlers.NewRouter(repo, s.Conf)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := http.ListenAndServe(s.Conf.Address, r); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	go storage.RunSaver(saver, repo, s.Conf.StoreInterval)

	<-done
}
