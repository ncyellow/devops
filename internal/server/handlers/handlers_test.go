package handlers

import (
	"github.com/ncyellow/devops/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGaugeHandler(t *testing.T) {
	type want struct {
		statusCode int
		nameMetric string
		typeMetric string
		value      float64
	}
	tests := []struct {
		name        string
		request     string
		contentType string
		want        want
	}{
		{
			name:        "add counter metric with correct data",
			request:     "/update/counter/testCounter/100",
			contentType: "text/plain",
			want: want{
				statusCode: http.StatusOK,
				nameMetric: "testCounter",
				typeMetric: "counter",
				value:      100,
			},
		},
		{
			name:        "add counter metric without id",
			request:     "/update/counter/",
			contentType: "text/plain",
			want: want{
				statusCode: http.StatusNotFound,
				nameMetric: "testCounter",
				typeMetric: "counter",
				value:      0,
			},
		},
		{
			name:        "counter invalid value",
			request:     "/update/counter/testCounter/none",
			contentType: "text/plain",
			want: want{
				statusCode: http.StatusBadRequest,
				nameMetric: "testCounter",
				typeMetric: "counter",
				value:      0,
			},
		},
		{
			name:        "add gauge metric with correct data",
			request:     "/update/gauge/testGauge/100",
			contentType: "text/plain",
			want: want{
				statusCode: http.StatusOK,
				nameMetric: "testGauge",
				typeMetric: "gauge",
				value:      100,
			},
		},
		{
			name:        "add gauge metric without id",
			request:     "/update/gauge/",
			contentType: "text/plain",
			want: want{
				statusCode: http.StatusNotFound,
				nameMetric: "testGauge",
				typeMetric: "gauge",
				value:      0,
			},
		},
		{
			name:        "gauge invalid value",
			request:     "/update/gauge/testGauge/none",
			contentType: "text/plain",
			want: want{
				statusCode: http.StatusBadRequest,
				nameMetric: "testGauge",
				typeMetric: "gauge",
				value:      0,
			},
		},
		{
			name:        "invalid update type",
			request:     "/update/unknown/testCounter/100",
			contentType: "text/plain",
			want: want{
				statusCode: http.StatusNotImplemented,
				nameMetric: "testCounter",
				typeMetric: "unknown",
				value:      0,
			},
		},
	}

	repo := storage.NewRepository()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			request.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(Handler(repo))
			h.ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()

			//! Проверяем что код ответа корректный
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			//! Если код ответа ок значит, то значит надо проверить корректно ли записалось значения метрики
			if result.StatusCode == http.StatusOK {
				if tt.want.typeMetric == "gauge" {
					val, _ := repo.Gauge(tt.want.nameMetric)
					assert.Equal(t, tt.want.value, val, tt.name)
				} else {
					val, _ := repo.Counter(tt.want.nameMetric)
					assert.Equal(t, int64(tt.want.value), val, tt.name)
				}
			}
		})
	}
}
