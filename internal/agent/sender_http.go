package agent

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ncyellow/devops/internal/agent/config"
	"github.com/ncyellow/devops/internal/crypto/rsa"
	"github.com/ncyellow/devops/internal/repository"
	"github.com/rs/zerolog/log"
)

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
