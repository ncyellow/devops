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
	"github.com/ncyellow/devops/internal/repository"
	"github.com/rs/zerolog/log"
)

// RunSender запускает цикл по обработке таймера отправки метрик из канала out на сервер
func RunSender(ctx context.Context, conf *config.Config, out <-chan []repository.Metrics, wg *sync.WaitGroup) {

	repo := repository.NewRepository(conf.GeneralCfg())
	encoder, err := rsa.NewEncoder(conf.CryptoKey)
	if err != nil {
		log.Info().Err(err)
	}
	url := fmt.Sprintf("http://%s/updates/", conf.Address)
	urlSingle := fmt.Sprintf("http://%s/update/", conf.Address)

	tickerReport := time.NewTicker(conf.ReportInterval.Duration)
	defer tickerReport.Stop()

	for {
		select {
		case <-tickerReport.C:
			// Две отправки для совместимости со старой версией, по старому протоколу
			SendMetrics(repo.ToMetrics(), urlSingle)
			// По новой через Batch
			SendMetricsBatch(repo.ToMetrics(), url, encoder)
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

// SendMetricsBatch отправляет все метрики одной пачкой на указанный url
func SendMetricsBatch(dataSource []repository.Metrics, url string, encoder *rsa.Encoder) {
	// Если метрик данных нет сразу на выход
	if len(dataSource) == 0 {
		return
	}

	buf, err := json.Marshal(dataSource)
	if err != nil {
		log.Fatal().Err(err)
	}
	if encoder != nil {
		buf, err = encoder.Encode(buf)
		if err != nil {
			log.Info().Msgf("проблемы с шифрованием отправлять не будем. %s", err.Error())
			return
		}
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
