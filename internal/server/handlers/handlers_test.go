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
		value      float64
	}
	tests := []struct {
		name        string
		request     string
		contentType string
		want        want
	}{
		{
			name:        "add gauge metric with correct data",
			request:     "/update/gauge/Alloc/1.00001",
			contentType: "text/plain",
			want: want{
				statusCode: http.StatusOK,
				nameMetric: "Alloc",
				value:      1.00001,
			},
		},
		{
			name:        "wrong contentType",
			request:     "/update/gauge/Alloc/1.00001",
			contentType: "application/json",
			want: want{
				statusCode: http.StatusInternalServerError,
				nameMetric: "Alloc",
				value:      1.00001,
			},
		},
		{
			name:        "add gauge metric with wrong name",
			request:     "/update/gauge/IncorrectName/1.021",
			contentType: "text/plain",
			want: want{
				statusCode: http.StatusInternalServerError,
				nameMetric: "IncorrectName",
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

			//! Проверяем что код ответа корректный
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			//! Если код ответа ок значит, то значит надо проверить корректно ли записалось значения метрики
			if result.StatusCode == http.StatusOK {
				val, _ := repo.Gauge(tt.want.nameMetric)
				assert.Equal(t, tt.want.value, val, tt.name)
			}
		})
	}
}
