package api

import (
	"context"
	"testing"

	pb "github.com/ncyellow/devops/internal/grpc/proto"
	"github.com/ncyellow/devops/internal/repository"
	"github.com/ncyellow/devops/internal/server/config"
	"github.com/ncyellow/devops/internal/server/storage"
	"github.com/stretchr/testify/assert"
)

func TestMetricsServer_AddMetric(t *testing.T) {
	conf := config.Config{}
	repo := repository.NewRepository(conf.GeneralCfg())
	store, err := storage.CreateStorage(&conf, repo)
	assert.NoError(t, err)

	server := NewMetricServer(repo, &conf, store)
	// Добавляем метрики
	response, err := server.AddMetric(context.Background(), &pb.AddMetricRequest{
		Counters: []*pb.CounterMetric{
			{
				Name:  "testCounter",
				Value: 100,
			},
		},
		Gauges: []*pb.GaugeMetric{
			{
				Name:  "testGauge",
				Value: 150,
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, response.Error, "")

	// Проверяем Get Counter метрики
	getResponse, err := server.GetMetric(context.Background(), &pb.GetMetricRequest{
		Name: "testCounter",
		Type: pb.Type_Counter,
	})
	assert.NoError(t, err)
	assert.NotNil(t, getResponse.Counter)
	assert.Nil(t, getResponse.Gauge)
	assert.Equal(t, getResponse.Counter.Value, int64(100))

	// Проверяем Get Gauge метрики
	getResponse, err = server.GetMetric(context.Background(), &pb.GetMetricRequest{
		Name: "testGauge",
		Type: pb.Type_Gauge,
	})
	assert.NoError(t, err)
	assert.NotNil(t, getResponse.Gauge)
	assert.Nil(t, getResponse.Counter)
	assert.Equal(t, getResponse.Gauge.Value, float64(150))
}
