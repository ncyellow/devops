package handlers

import (
	"fmt"
	"github.com/ncyellow/devops/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

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
		want        want
	}{
		{
			name:        "add counter metric with correct data",
			request:     "/update/counter/testCounter/100",
			requestType: "POST",
			want: want{
				statusCode: http.StatusOK,
				body:       "ok",
			},
		},
		{
			name:        "add counter metric without id",
			request:     "/update/counter/",
			requestType: "POST",
			want: want{
				statusCode: http.StatusNotFound,
				body:       "404 page not found\n",
			},
		},
		{
			name:        "counter invalid value",
			request:     "/update/counter/testCounter/none",
			requestType: "POST",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "incorrect metric value",
			},
		},
		{
			name:        "add gauge metric with correct data",
			request:     "/update/gauge/testGauge/100",
			requestType: "POST",
			want: want{
				statusCode: http.StatusOK,
				body:       "ok",
			},
		},
		{
			name:        "add gauge metric without id",
			request:     "/update/gauge/",
			requestType: "POST",
			want: want{
				statusCode: http.StatusNotFound,
				body:       "404 page not found\n",
			},
		},
		{
			name:        "gauge invalid value",
			request:     "/update/gauge/testGauge/none",
			requestType: "POST",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "incorrect metric value",
			},
		},
		{
			name:        "invalid update type",
			request:     "/update/unknown/testCounter/100",
			requestType: "POST",
			want: want{
				statusCode: http.StatusNotImplemented,
				body:       "incorrect metric type",
			},
		},
		{
			name:        "get correct counter value",
			request:     "/value/counter/testCounter",
			requestType: "GET",
			want: want{
				statusCode: http.StatusOK,
				body:       "100",
			},
		},
		{
			name:        "get unknown counter value",
			request:     "/value/counter/unknownCounter",
			requestType: "GET",
			want: want{
				statusCode: http.StatusNotFound,
				body:       "",
			},
		},
		{
			name:        "get correct gauge value",
			request:     "/value/gauge/testGauge",
			requestType: "GET",
			want: want{
				statusCode: http.StatusOK,
				body:       fmt.Sprintf("%.3f", 100.0),
			},
		},
		{
			name:        "get unknown gauge value",
			request:     "/value/gauge/unknownGauge",
			requestType: "GET",
			want: want{
				statusCode: http.StatusNotFound,
				body:       "",
			},
		},
		{
			name:        "list all metrics",
			request:     "/",
			requestType: "GET",
			want: want{
				statusCode: http.StatusOK,
				body:       "this is all metrics",
			},
		},
	}

	for _, tt := range tests {
		resp, body := testRequest(t, ts, tt.requestType, tt.request)
		assert.Equal(t, tt.want.statusCode, resp.StatusCode, tt.name)
		assert.Equal(t, tt.want.body, body, tt.name)
	}

}
