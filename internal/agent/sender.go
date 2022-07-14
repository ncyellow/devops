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
	"github.com/ncyellow/devops/internal/server/repository"
	"github.com/rs/zerolog/log"
)

// RunSender запускает цикл по обработке таймеров и ожидания сигналов от ОС
func RunSender(ctx context.Context, conf *config.Config, out <-chan []repository.Metrics, wg *sync.WaitGroup) {

	repo := repository.NewRepository(conf.GeneralCfg())
	url := fmt.Sprintf("http://%s/updates/", conf.Address)
	urlSingle := fmt.Sprintf("http://%s/update/", conf.Address)

	tickerReport := time.NewTicker(conf.ReportInterval)
	defer tickerReport.Stop()

	for {
		select {
		case <-tickerReport.C:
			SendMetrics(repo.ToMetrics(), urlSingle)
			SendMetricsBatch(repo.ToMetrics(), url)
		case metrics := <-out:
			//! Корректный выход без ошибок по указанным сигналам
			// обработка входящих метрик
			fmt.Println("got metrics to sender")
			for _, metric := range metrics {
				repo.UpdateMetric(metric)
			}
		case <-ctx.Done():
			wg.Done()
			return
		}
	}
}

// SendMetricsBatch отправляет все метрики одной пачкой на указанный url
func SendMetricsBatch(dataSource []repository.Metrics, url string) {
	if len(dataSource) == 0 {
		return
	}

	buf, err := json.Marshal(dataSource)
	if err != nil {
		log.Fatal().Err(err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		log.Info().Msgf("%s", err.Error())
		return
	}
	resp.Body.Close()
}

// SendMetrics отправляет метрики на указанный url
func SendMetrics(dataSource []repository.Metrics, url string) {
	for _, metric := range dataSource {
		buf, err := json.Marshal(metric)
		if err != nil {
			log.Fatal().Err(err)
		}
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(buf))
		if err != nil {
			log.Info().Msgf("%s", err.Error())
			continue
		}
		resp.Body.Close()
	}
}
