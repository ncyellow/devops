package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// Проверяем основные варианты. Включена / выключена поддержка gzip
func TestEncoderGZIP(t *testing.T) {

	tests := map[string]struct {
		request          func() *http.Request
		expectedResponse string
	}{
		"Request without gzip support": {
			func() *http.Request {
				req, _ := http.NewRequest("GET", "/", nil)
				return req
			},
			"ok",
		},

		"Request with gzip support": {
			func() *http.Request {
				req, _ := http.NewRequest("GET", "/", nil)
				req.Header.Add("Accept-Encoding", "gzip")
				return req
			},
			"\x1f\x8b\b\x00\x00\x00\x00\x00\x04\xff\x00\x02\x00\xfd\xffok\x01\x00\x00\xff\xffG\xdd\xdcy\x02\x00\x00\x00",
		},
	}

	r := chi.NewRouter()
	r.Use(EncoderGZIP)
	for _, test := range tests {
		w := httptest.NewRecorder()

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			response := "ok"
			w.Write([]byte(response))
		})
		r.ServeHTTP(w, test.request())

		assert.Equal(t, w.Body.String(), test.expectedResponse)
	}
}
