package server

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/ncyellow/devops/internal/grpc/api"
	pb "github.com/ncyellow/devops/internal/grpc/proto"
	"github.com/ncyellow/devops/internal/repository"
	"github.com/ncyellow/devops/internal/server/config"
	"github.com/ncyellow/devops/internal/server/storage"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// GRPCServer структура сервера
type GRPCServer struct {
	Conf *config.Config
}

// RunServer блокирующая функция запуска сервера.
// После запуска встает в ожидание os.Interrupt, syscall.SIGINT, syscall.SIGTERM
func (s *GRPCServer) RunServer() {
	repo := repository.NewRepository(s.Conf.GeneralCfg())

	saver, err := storage.CreateStorage(s.Conf, repo)
	if err != nil {
		log.Info().Msg("cant create NewPgStorage")
	}
	defer saver.Close()
	// Поднимаем текущие данные по метриками
	saver.Load()

	listen, err := net.Listen("tcp", s.Conf.GRPCAddress)
	if err != nil {
		log.Fatal().Err(err)
	}

	grpcServer := grpc.NewServer()
	// регистрируем сервис
	pb.RegisterMetricsServer(grpcServer, api.NewMetricServer(repo, s.Conf, saver))

	defer func() {
		// гасим сервер через GracefulStop
		grpcServer.GracefulStop()
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		if err := grpcServer.Serve(listen); err != nil {
			log.Error().Err(err)
		}
	}()

	go storage.RunStorageSaver(saver, s.Conf.StoreInterval.Duration)

	<-done
	log.Info().Msg("Server Shutdown gracefully")

}
