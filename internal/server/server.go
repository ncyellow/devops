// Package server содержит кодовую базу сервера.
// Конкретно данный файл содержит основной код запуска сервера:
// Реализация определяется http grpc по конфигурации
// server := server.CreateServer(&cfg)
// server.RunServer()
package server

import "github.com/ncyellow/devops/internal/server/config"

// Server интерфейс сервера
type Server interface {
	// RunServer - запуск в синхронном режиме
	RunServer()
}

// CreateServer - Основная factory функция создания сервера. По конфигурации выбирает реализацию
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
