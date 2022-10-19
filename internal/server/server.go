// Package server содержит кодовую базу сервера.
// Конкретно данный файл содержит основной код запуска сервера:
// server := server.Server{Conf: &cfg}
// server.RunServer()
package server

import "github.com/ncyellow/devops/internal/server/config"

type Server interface {
	RunServer()
}

func CreateServer(conf *config.Config) Server {
	if conf.GRPCAddress != "" {
		return &GRPCServer{
			Conf: conf,
		}
	}
	return &HTTPServer{
		Conf: conf,
	}
}
