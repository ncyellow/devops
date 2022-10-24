package agent

import (
	"context"

	"github.com/ncyellow/devops/internal/agent/config"
	"github.com/ncyellow/devops/internal/grpc/proto"
	"github.com/ncyellow/devops/internal/repository"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// GRPCSender структура для отправки на сервер по grpc
type GRPCSender struct {
	conf   *config.Config
	conn   *grpc.ClientConn
	client proto.MetricsClient
}

// SendMetricsBatch отправляет все метрики одной пачкой на указанный url
func (g *GRPCSender) SendMetricsBatch(dataSource []repository.Metrics) {
	// Если метрик данных нет сразу на выход
	if len(dataSource) == 0 {
		return
	}

	var counters []*proto.CounterMetric
	var gauges []*proto.GaugeMetric
	for _, metric := range dataSource {
		// делаем так потому что hash опциональное поле и мы будем передавать указатель и потому отдельная переменная
		var hash string
		if metric.Hash != "" {
			hash = metric.Hash
		}

		switch metric.MType {
		case repository.Counter:
			counters = append(counters, &proto.CounterMetric{
				Name:  metric.ID,
				Value: *metric.Delta,
				Hash:  &hash,
			})
		case repository.Gauge:
			gauges = append(gauges, &proto.GaugeMetric{
				Name:  metric.ID,
				Value: *metric.Value,
				Hash:  &hash,
			})
		}
	}

	resp, err := g.client.AddMetric(context.Background(), &proto.AddMetricRequest{
		Counters: counters,
		Gauges:   gauges,
	})
	if err != nil {
		log.Info().Msgf("%s", err.Error())
		return
	}
	if resp.Error != "" {
		log.Info().Msg(resp.Error)
	}
}

// SendMetrics отправляет все метрики одной пачкой на указанный url,
// это мы делаем только http для совместимости со старыми автотестами
func (g *GRPCSender) SendMetrics(dataSource []repository.Metrics) {
	// Если метрик данных нет сразу на выход
	if len(dataSource) == 0 {
		return
	}
}

// Close общая функция очистки ресурсов. Для http не требуется
func (g *GRPCSender) Close() {
	g.conn.Close()
}
