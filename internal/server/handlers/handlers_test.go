package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ncyellow/devops/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, contentType string, reqBody []byte) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, bytes.NewBuffer(reqBody))
	require.NoError(t, err)

	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

// TestRouter тесты по запросам к различным url
func TestRouter(t *testing.T) {
	repo := storage.NewRepository()
	r := NewRouter(repo)

	ts := httptest.NewServer(r)
	defer ts.Close()

	type want struct {
		statusCode int
		body       string
	}
	tests := []struct {
		name        string
		request     string
		requestType string
		contentType string
		body        []byte
		want        want
	}{
		{
			name:        "add counter metric with correct data",
			request:     "/update/counter/testCounter/100",
			requestType: "POST",
			contentType: "text/plain",
			body:        nil,
			want: want{
				statusCode: http.StatusOK,
				body:       "ok",
			},
		},
		{
			name:        "add counter metric without id",
			request:     "/update/counter/",
			requestType: "POST",
			contentType: "text/plain",
			body:        nil,
			want: want{
				statusCode: http.StatusNotFound,
				body:       "404 page not found\n",
			},
		},
		{
			name:        "counter invalid value",
			request:     "/update/counter/testCounter/none",
			requestType: "POST",
			contentType: "text/plain",
			body:        nil,
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "incorrect metric value",
			},
		},
		{
			name:        "add gauge metric with correct data",
			request:     "/update/gauge/testGauge/100",
			requestType: "POST",
			contentType: "text/plain",
			body:        nil,
			want: want{
				statusCode: http.StatusOK,
				body:       "ok",
			},
		},
		{
			name:        "add gauge metric without id",
			request:     "/update/gauge/",
			requestType: "POST",
			contentType: "text/plain",
			body:        nil,
			want: want{
				statusCode: http.StatusNotFound,
				body:       "404 page not found\n",
			},
		},
		{
			name:        "gauge invalid value",
			request:     "/update/gauge/testGauge/none",
			requestType: "POST",
			contentType: "text/plain",
			body:        nil,
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "incorrect metric value",
			},
		},
		{
			name:        "invalid update type",
			request:     "/update/unknown/testCounter/100",
			requestType: "POST",
			contentType: "text/plain",
			body:        nil,
			want: want{
				statusCode: http.StatusNotImplemented,
				body:       "incorrect metric type",
			},
		},
		{
			name:        "get correct counter value",
			request:     "/value/counter/testCounter",
			requestType: "GET",
			contentType: "text/plain",
			body:        nil,
			want: want{
				statusCode: http.StatusOK,
				body:       "100",
			},
		},
		{
			name:        "get unknown counter value",
			request:     "/value/counter/unknownCounter",
			requestType: "GET",
			contentType: "text/plain",
			body:        nil,
			want: want{
				statusCode: http.StatusNotFound,
				body:       "",
			},
		},
		{
			name:        "get correct gauge value",
			request:     "/value/gauge/testGauge",
			requestType: "GET",
			contentType: "text/plain",
			body:        nil,
			want: want{
				statusCode: http.StatusOK,
				body:       fmt.Sprintf("%.3f", 100.0),
			},
		},
		{
			name:        "get unknown gauge value",
			request:     "/value/gauge/unknownGauge",
			requestType: "GET",
			contentType: "text/plain",
			body:        nil,
			want: want{
				statusCode: http.StatusNotFound,
				body:       "",
			},
		},
		{
			name:        "list all metrics",
			request:     "/",
			requestType: "GET",
			contentType: "text/plain",
			body:        nil,
			want: want{
				statusCode: http.StatusOK,
				body: `
	<html>
	<body>
	<h1>All metrics</h1>
	<h3>gauges</h3>
	<ul>
	  <li>testGauge : 100.000</li>

	</ul>
	<h3>counters</h3>
	<ul>
	  <li>testCounter : 100</li>

	</ul>
	</body>
	</html>`,
			},
		},
		{
			name:        "set gauge with json",
			request:     "/update/",
			requestType: "POST",
			contentType: "application/json",
			body:        []byte(`{"id":"jsonGauge","type":"gauge","value": 111}`),
			want: want{
				statusCode: http.StatusOK,
				body:       "ok",
			},
		},
		{
			name:        "set counter with json",
			request:     "/update/",
			requestType: "POST",
			contentType: "application/json",
			body:        []byte(`{"id":"jsonCounter","type":"counter","delta": 123}`),
			want: want{
				statusCode: http.StatusOK,
				body:       "ok",
			},
		},
		{
			name:        "test get gauge with json",
			request:     "/value/",
			requestType: "POST",
			contentType: "application/json",
			body:        []byte(`{"id":"jsonGauge","type":"gauge"}`),
			want: want{
				statusCode: http.StatusOK,
				body:       `{"id":"jsonGauge","type":"gauge","value":111}`,
			},
		},
		{
			name:        "test get counter with json",
			request:     "/value/",
			requestType: "POST",
			contentType: "application/json",
			body:        []byte(`{"id":"jsonCounter","type":"counter"}`),
			want: want{
				statusCode: http.StatusOK,
				body:       `{"id":"jsonCounter","type":"counter","delta":123}`,
			},
		},
	}

	for _, tt := range tests {
		resp, body := testRequest(t, ts, tt.requestType, tt.request, tt.contentType, tt.body)
		assert.Equal(t, tt.want.statusCode, resp.StatusCode, tt.name)
		assert.Equal(t, tt.want.body, body, tt.name)
		resp.Body.Close()
	}

}
