package agent

import (
	"context"
	"github.com/ncyellow/devops/internal/agent/config"
	"github.com/ncyellow/devops/internal/genconfig"
	"github.com/ncyellow/devops/internal/repository"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestCollector(t *testing.T) {
	conf := genconfig.GeneralConfig{
		Address:   "Test",
		SecretKey: "Test",
	}
	collector := Collector{
		Conf:   &conf,
		Source: &RuntimeSource{},
	}
	collector.Update()
	metrics := collector.ToMetrics()
	assert.Equal(t, len(metrics), 29)
}

func TestRunCollector(t *testing.T) {
	aconf := config.Config{
		GeneralConfig: genconfig.GeneralConfig{
			Address:   "Test",
			SecretKey: "Test",
		},
		PollInterval:   time.Second * 1,
		ReportInterval: time.Second * 1,
	}
	collector := Collector{
		Conf:   aconf.GeneralCfg(),
		Source: &RuntimeSource{},
	}

	metricChannel := make(chan []repository.Metrics, 1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go RunCollector(ctx, &aconf, &collector, metricChannel, &wg)

	go func() {
		wg.Wait()
		close(metricChannel)
	}()

	for {
		select {
		case metrics := <-metricChannel:
			assert.NotNil(t, metrics)
		case <-ctx.Done():
			return
		}
	}
}
