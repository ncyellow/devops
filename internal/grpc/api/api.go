// Package api - реализует апи grpc сервера.
package api

import (
	"context"

	"github.com/ncyellow/devops/internal/crypto/rsa"
	"github.com/ncyellow/devops/internal/grpc/proto"
	"github.com/ncyellow/devops/internal/hash"
	"github.com/ncyellow/devops/internal/repository"
	"github.com/ncyellow/devops/internal/server/config"
	"github.com/ncyellow/devops/internal/server/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MetricsServer Реализуем все базовые возможности grpc сервера
// - установить значение списка метрик
// - прочитать значение метрики
// - получить html
// - пинг
type MetricsServer struct {
	proto.UnimplementedMetricsServer
	conf    *config.Config
	repo    repository.Repository
	pStore  storage.PersistentStorage
	decoder *rsa.Decoder
}

func NewMetricServer(repo repository.Repository, conf *config.Config, pStore storage.PersistentStorage) *MetricsServer {
	return &MetricsServer{
		repo:   repo,
		conf:   conf,
		pStore: pStore,
	}
}

func (ms *MetricsServer) AddMetric(ctx context.Context, req *proto.AddMetricRequest) (*proto.AddMetricResponse, error) {
	var response proto.AddMetricResponse
	encodeFunc := hash.CreateEncodeFunc(ms.conf.SecretKey)
	counters := req.GetCounters()
	for _, metric := range counters {
		value := metric.GetValue()
		counter := repository.Metrics{
			ID:    metric.GetName(),
			MType: repository.Counter,
			Delta: &value,
		}
		if metric.Hash != nil {
			counter.Hash = *metric.Hash
		}

		ok := hash.CheckSign(ms.conf.SecretKey, counter.Hash, counter.CalcHash(encodeFunc))
		if !ok {
			return nil, status.Error(codes.InvalidArgument, "incorrect metric sign")
		}
		ms.repo.UpdateMetric(counter)
	}
	gauges := req.GetGauges()
	for _, metric := range gauges {
		delta := metric.GetValue()
		gauge := repository.Metrics{
			ID:    metric.GetName(),
			MType: repository.Gauge,
			Value: &delta,
		}
		if metric.Hash != nil {
			gauge.Hash = *metric.Hash
		}

		ok := hash.CheckSign(ms.conf.SecretKey, gauge.Hash, gauge.CalcHash(encodeFunc))
		if !ok {
			return nil, status.Error(codes.InvalidArgument, "incorrect metric sign")
		}
		ms.repo.UpdateMetric(gauge)
	}
	return &response, nil
}
func (ms *MetricsServer) GetMetric(ctx context.Context, req *proto.GetMetricRequest) (*proto.GetMetricResponse, error) {
	var response proto.GetMetricResponse
	switch req.GetType() {
	case proto.Type_Counter:
		val, ok := ms.repo.Counter(req.GetName())
		if !ok {
			return nil, status.Errorf(codes.NotFound, "not found")
		}
		response.Counter = &proto.CounterMetric{
			Name:  req.GetName(),
			Value: val,
		}
	case proto.Type_Gauge:
		val, ok := ms.repo.Gauge(req.GetName())
		if !ok {
			return nil, status.Errorf(codes.NotFound, "not found")
		}
		response.Gauge = &proto.GaugeMetric{
			Name:  req.GetName(),
			Value: val,
		}
	}
	return &response, nil
}

func (ms *MetricsServer) ListMetrics(context.Context, *proto.ListMetricsRequest) (*proto.ListMetricResponse, error) {
	var response proto.ListMetricResponse
	response.Html = repository.RenderHTML(ms.repo.ToMetrics())
	return &response, nil
}
func (ms *MetricsServer) Ping(context.Context, *proto.PingRequest) (*proto.PingResponse, error) {
	var response proto.PingResponse
	err := ms.pStore.Ping()
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &response, nil
}
