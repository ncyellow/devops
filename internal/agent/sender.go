package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

type HTTPSender struct {
	conf      *config.Config
	urlBatch  string
	urlSingle string
	encoder   *rsa.Encoder
}

// SendMetricsBatch отправляет все метрики одной пачкой на указанный url
func (s *HTTPSender) SendMetricsBatch(dataSource []repository.Metrics) {
	// Если метрик данных нет сразу на выход
	if len(dataSource) == 0 {
		return
	}

	buf, err := json.Marshal(dataSource)
	if err != nil {
		log.Fatal().Err(err)
	}
	if s.encoder != nil {
		buf, err = s.encoder.Encode(buf)
		if err != nil {
			log.Info().Msgf("проблемы с шифрованием отправлять не будем. %s", err.Error())
			return
		}
	}

	resp, err := http.Post(s.urlBatch, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		log.Info().Msgf("%s", err.Error())
		return
	}
	resp.Body.Close()
}

// SendMetrics отправляет метрики на указанный url
func (s *HTTPSender) SendMetrics(dataSource []repository.Metrics) {
	client := http.Client{Timeout: 100 * time.Millisecond}
	for _, metric := range dataSource {
		buf, err := json.Marshal(metric)
		if err != nil {
			log.Fatal().Err(err)
		}
		resp, err := client.Post(s.urlSingle, "application/json", bytes.NewBuffer(buf))
		if err != nil {
			log.Info().Msgf("%s", err.Error())
			continue
		}
		resp.Body.Close()
	}
}

func (s *HTTPSender) Close() {
	// Общая функция очистки ресурсов. для http не требуется
}

func CreateSender(conf *config.Config) Sender {
	if conf.Address != "" {
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
	conn, err := grpc.Dial(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
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

type GRPCSender struct {
	conf   *config.Config
	conn   *grpc.ClientConn
	client pb.MetricsClient
}

// SendMetricsBatch отправляет все метрики одной пачкой на указанный url
func (g *GRPCSender) SendMetricsBatch(dataSource []repository.Metrics) {
	// Если метрик данных нет сразу на выход
	if len(dataSource) == 0 {
		return
	}

	var counters []*pb.CounterMetric
	var gauges []*pb.GaugeMetric
	for _, metric := range dataSource {
		switch metric.MType {
		case repository.Counter:
			counters = append(counters, &pb.CounterMetric{
				Name:  metric.ID,
				Value: *metric.Delta,
			})
		case repository.Gauge:
			gauges = append(gauges, &pb.GaugeMetric{
				Name:  metric.ID,
				Value: *metric.Value,
			})
		}
	}

	resp, err := g.client.AddMetric(context.Background(), &pb.AddMetricRequest{
		Counters: counters,
		Gauges:   gauges,
	})
	if err != nil {
		log.Info().Msgf("%s", err.Error())
	}
	if resp.Error != "" {
		log.Info().Msg(resp.Error)
	}
}

// SendMetrics отправляет все метрики одной пачкой на указанный url
func (g *GRPCSender) SendMetrics(dataSource []repository.Metrics) {
	// Если метрик данных нет сразу на выход
	if len(dataSource) == 0 {
		return
	}
}

func (g *GRPCSender) Close() {
	// Общая функция очистки ресурсов. для http не требуется
	g.conn.Close()
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
