package server

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ncyellow/devops/internal/server/config"
	"github.com/ncyellow/devops/internal/server/handlers"
	"github.com/ncyellow/devops/internal/server/repository"
	"github.com/ncyellow/devops/internal/server/storage"
	"github.com/rs/zerolog/log"
)

type Server struct {
	Conf *config.Config
}

func (s Server) RunServer() {
	repo := repository.NewRepository(s.Conf)

	saver, err := storage.CreateStorage(s.Conf, repo)
	if err != nil {
		log.Info().Msg("cant create NewSaver")
	}
	defer saver.Close()

	saver.Load()

	r := handlers.NewRouter(repo, s.Conf, saver)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := http.ListenAndServe(s.Conf.Address, r); err != nil && err != http.ErrServerClosed {
			log.Error().Msgf("listen: %s", err)
		}
	}()

	go storage.RunSaver(saver, s.Conf.StoreInterval)

	<-done
}
