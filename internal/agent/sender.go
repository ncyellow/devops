package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ncyellow/devops/internal/agent/config"
	"github.com/ncyellow/devops/internal/crypto/rsa"
	pb "github.com/ncyellow/devops/internal/grpc/proto"
	"github.com/ncyellow/devops/internal/repository"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Sender interface {
	SendMetricsBatch(dataSource []repository.Metrics)
	SendMetrics(dataSource []repository.Metrics)
	Close()
}

func CreateSender(conf *config.Config) Sender {
	// По дефолту у нас http, только если задан GRPCAddress entrypoint, мы переходим на grpc
	if conf.GRPCAddress == "" {
		encoder, err := rsa.NewEncoder(conf.CryptoKey)
		if err != nil {
			log.Info().Err(err)
		}
		return &HTTPSender{
			conf:      conf,
			urlBatch:  fmt.Sprintf("http://%s/updates/", conf.Address),
			urlSingle: fmt.Sprintf("http://%s/update/", conf.Address),
			encoder:   encoder,
		}
	}

	// устанавливаем соединение с сервером
	conn, err := grpc.Dial(conf.GRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Err(err)
	}
	// получаем переменную интерфейсного типа UsersClient,
	// через которую будем отправлять сообщения
	client := pb.NewMetricsClient(conn)
	return &GRPCSender{
		conf:   conf,
		conn:   conn,
		client: client,
	}
}

// RunSender запускает цикл по обработке таймера отправки метрик из канала out на сервер
func RunSender(ctx context.Context, conf *config.Config, out <-chan []repository.Metrics, wg *sync.WaitGroup) {

	repo := repository.NewRepository(conf.GeneralCfg())

	sender := CreateSender(conf)
	defer sender.Close()

	tickerReport := time.NewTicker(conf.ReportInterval.Duration)
	defer tickerReport.Stop()

	for {
		select {
		case <-tickerReport.C:
			// Две отправки для совместимости со старой версией, по старому протоколу
			sender.SendMetrics(repo.ToMetrics())
			// По новой через Batch
			sender.SendMetricsBatch(repo.ToMetrics())
		case metrics, ok := <-out:
			if !ok {
				wg.Done()
				return
			}
			for _, metric := range metrics {
				repo.UpdateMetric(metric)
			}
		case <-ctx.Done():
			wg.Done()
			return
		}
	}
}
