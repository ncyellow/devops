package api

import (
	"context"
	"testing"

	"github.com/ncyellow/devops/internal/grpc/proto"
	"github.com/ncyellow/devops/internal/repository"
	"github.com/ncyellow/devops/internal/server/config"
	"github.com/ncyellow/devops/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMetricsServer(t *testing.T) {
	conf := config.Config{}
	repo := repository.NewRepository(conf.GeneralCfg())
	store, err := storage.CreateStorage(&conf, repo)
	assert.NoError(t, err)

	server := NewMetricServer(repo, &conf, store)
	// Добавляем метрики
	response, err := server.AddMetric(context.Background(), &proto.AddMetricRequest{
		Counters: []*proto.CounterMetric{
			{
				Name:  "testCounter",
				Value: 100,
			},
		},
		Gauges: []*proto.GaugeMetric{
			{
				Name:  "testGauge",
				Value: 150,
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, response.Error, "")

	// Проверяем Get Counter метрики
	getResponse, err := server.GetMetric(context.Background(), &proto.GetMetricRequest{
		Name: "testCounter",
		Type: proto.Type_Counter,
	})
	assert.NoError(t, err)
	assert.NotNil(t, getResponse.Counter)
	assert.Nil(t, getResponse.Gauge)
	assert.Equal(t, getResponse.Counter.Value, int64(100))

	// Проверяем Get Gauge метрики
	getResponse, err = server.GetMetric(context.Background(), &proto.GetMetricRequest{
		Name: "testGauge",
		Type: proto.Type_Gauge,
	})
	assert.NoError(t, err)
	assert.NotNil(t, getResponse.Gauge)
	assert.Nil(t, getResponse.Counter)
	assert.Equal(t, getResponse.Gauge.Value, float64(150))

	pingResponse, err := server.Ping(context.Background(), &proto.PingRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, pingResponse.Error)

	listResponse, err := server.ListMetrics(context.Background(), &proto.ListMetricsRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listResponse.Html)
	assert.Equal(t, listResponse.Html, `
	<html>
	<body>
	<h1>All metrics</h1>
	<h3>gauges</h3>
	<ul>
	  <li>testGauge : 150.000</li>

	</ul>
	<h3>counters</h3>
	<ul>
	  <li>testCounter : 100</li>

	</ul>
	</body>
	</html>`)
}

func TestMetricsServer_GetMetric(t *testing.T) {
	//! Отдельный тест проверок запросов не корректных значений

	conf := config.Config{}
	repo := repository.NewRepository(conf.GeneralCfg())
	store, err := storage.CreateStorage(&conf, repo)
	assert.NoError(t, err)

	server := NewMetricServer(repo, &conf, store)

	metricTypes := []proto.Type{
		proto.Type_Counter,
		proto.Type_Gauge,
	}

	// Проверяем что для всех типов метрик запрос неизвестных метрик выдает ошибку
	for _, mType := range metricTypes {
		// Проверяем Get Counter - будет ошибка если метрика на сервере не найдена
		getResponse, err := server.GetMetric(context.Background(), &proto.GetMetricRequest{
			Name: "testCounter",
			Type: mType,
		})
		assert.Error(t, err)
		assert.Nil(t, getResponse)
		s, ok := status.FromError(err)
		assert.True(t, ok, true)
		assert.Equal(t, s.Code(), codes.NotFound)
	}
}
