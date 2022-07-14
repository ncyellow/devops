package agent

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ncyellow/devops/internal/agent/config"
	"github.com/ncyellow/devops/internal/server/repository"
)

type Agent struct {
	Conf    *config.Config
	metrics RuntimeMetrics
}

// Run запускает цикл по обработке таймеров и ожидания сигналов от ОС
func (collector *Agent) Run() error {
	ctx, cancel := context.WithCancel(context.Background())

	var metricChannel chan []repository.Metrics
	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go RunCollector(ctx, collector.Conf, metricChannel, &wg)
	wg.Add(1)
	go RunSender(ctx, collector.Conf, metricChannel, &wg)

	<-signalChanel
	cancel()
	fmt.Println("ok")
	return nil
}
