// Package server содержит кодовую базу сервера.
// Конкретно данный файл содержит основной код запуска сервера:
// server := server.Server{Conf: &cfg}
// server.RunServer()
package server

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ncyellow/devops/internal/repository"
	"github.com/ncyellow/devops/internal/server/config"
	"github.com/ncyellow/devops/internal/server/handlers"
	"github.com/ncyellow/devops/internal/server/storage"
	"github.com/rs/zerolog/log"
)

// Server структура сервера
type Server struct {
	Conf *config.Config
}

// RunServer блокирующая функция запуска сервера.
// После запуска встает в ожидание os.Interrupt, syscall.SIGINT, syscall.SIGTERM
func (s Server) RunServer() {
	repo := repository.NewRepository(s.Conf.GeneralCfg())

	saver, err := storage.CreateStorage(s.Conf, repo)
	if err != nil {
		log.Info().Msg("cant create NewPgStorage")
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

	go storage.RunStorageSaver(saver, s.Conf.StoreInterval.Duration)

	<-done
}
