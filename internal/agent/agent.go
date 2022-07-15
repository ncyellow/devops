package agent

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ncyellow/devops/internal/agent/config"
	"github.com/ncyellow/devops/internal/server/repository"
)

type Agent struct {
	Conf *config.Config
}

// Run запускает цикл по обработке таймеров и ожидания сигналов от ОС
func (collector *Agent) Run() error {

	// Контекст для корректно завершения все горутин
	ctx, cancel := context.WithCancel(context.Background())

	// Канал по которому метрики откуда sender модуль будет получать метрики и отправлять на сервер
	metricChannel := make(chan []repository.Metrics, 1)

	// Канал обработки сигналов ОС
	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	wg := sync.WaitGroup{}

	// Коллектор по сбору runtime метрик
	runtimeCol := &Collector{
		Conf:   collector.Conf.GeneralCfg(),
		Source: &RuntimeSource{},
	}
	wg.Add(1)
	go RunCollector(ctx, collector.Conf, runtimeCol, metricChannel, &wg)

	// Коллектор по сбору psutil метрик
	psUtilCol := &Collector{
		Conf:   collector.Conf.GeneralCfg(),
		Source: NewPSUtilSource(),
	}
	wg.Add(1)
	go RunCollector(ctx, collector.Conf, psUtilCol, metricChannel, &wg)

	wg.Add(1)
	go RunSender(ctx, collector.Conf, metricChannel, &wg)

	<-signalChanel
	cancel()
	wg.Wait()
	close(metricChannel)
	return nil
}
