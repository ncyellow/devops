// Package server содержит кодовую базу сервера.
// Конкретно данный файл содержит основной код запуска сервера:
// server := server.Server{Conf: &cfg}
// server.RunServer()
package server

import (
	"context"
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

// HTTPServer структура сервера
type HTTPServer struct {
	Conf *config.Config
}

// RunServer блокирующая функция запуска сервера.
// После запуска встает в ожидание os.Interrupt, syscall.SIGINT, syscall.SIGTERM
func (s *HTTPServer) RunServer() {
	repo := repository.NewRepository(s.Conf.GeneralCfg())

	saver, err := storage.CreateStorage(s.Conf, repo)
	if err != nil {
		log.Info().Msg("cant create NewPgStorage")
	}
	defer saver.Close()
	// Поднимаем текущие данные по метриками
	saver.Load()

	srv := http.Server{
		Addr:    s.Conf.Address,
		Handler: handlers.NewRouter(repo, s.Conf, saver),
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	idleConnsClosed := make(chan struct{})

	go func() {
		// ждем прерывание
		<-done
		// гасим сервер
		if err := srv.Shutdown(context.Background()); err != nil {
			// ошибки закрытия Listener
			log.Info().Msgf("HTTP server Shutdown: %v", err)
		}
		// сообщаем основному потоку,
		// что все сетевые соединения обработаны и закрыты
		close(idleConnsClosed)
	}()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Msgf("listen: %s", err)
		}
	}()

	go storage.RunStorageSaver(saver, s.Conf.StoreInterval.Duration)

	<-idleConnsClosed
	log.Info().Msg("Server Shutdown gracefully")

}
