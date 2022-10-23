package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ncyellow/devops/internal/genconfig"
	"github.com/ncyellow/devops/internal/repository"
	"github.com/ncyellow/devops/internal/server/config"
	"github.com/ncyellow/devops/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type want struct {
	statusCode int
	body       string
}
type tests struct {
	name        string
	request     string
	requestType string
	contentType string
	body        []byte
	want        want
}

type HandlersSuite struct {
	suite.Suite
	ts *httptest.Server
}

// SetupSuite перед началом теста стартуем новый сервер httptest.Server делаем так, чтобы тестировать каждый
// handler отдельно и не сливать все тесты в один
func (suite *HandlersSuite) SetupTest() {
	conf := config.Config{
		GeneralConfig: genconfig.GeneralConfig{
			CryptoKey: "unknown_rsa",
		},
	}
	repo := repository.NewRepository(conf.GeneralCfg())
	//! Это пустой вариант хранилища без состояние. Ошибок нет
	pStore, _ := storage.NewFakeStorage()

	r := NewRouter(repo, &conf, pStore)
	suite.ts = httptest.NewServer(r)
}

// TearDownSuite после теста отключаем сервер
func (suite *HandlersSuite) TearDownTest() {
	suite.ts.Close()
}

func TestHandlersSuite(t *testing.T) {
	suite.Run(t, new(HandlersSuite))
}

func (suite *HandlersSuite) runTableTests(testList []tests) {
	for _, tt := range testList {
		resp, body := runTestRequest(suite.T(), suite.ts, tt.requestType, tt.request, tt.contentType, tt.body)
		assert.Equal(suite.T(), tt.want.statusCode, resp.StatusCode, tt.name)
		assert.Equal(suite.T(), tt.want.body, body, tt.name)
		resp.Body.Close()
	}
}

func runTestRequest(t *testing.T, ts *httptest.Server, method, path string, contentType string, reqBody []byte) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, bytes.NewBuffer(reqBody))
	require.NoError(t, err)

	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

//TestListHandler тестируем ListHandler
func (suite *HandlersSuite) TestListHandler() {
	testData := []tests{
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
			name:        "add gauge metric with incorrect value",
			request:     "/update/gauge/testGauge/10sa0",
			requestType: "POST",
			contentType: "text/plain",
			body:        nil,
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "incorrect metric value",
			},
		},
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
			name:        "add counter metric with incorrect value",
			request:     "/update/counter/testCounter/1dd00",
			requestType: "POST",
			contentType: "text/plain",
			body:        nil,
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "incorrect metric value",
			},
		},
		{
			name:        "add metric with incorrect type",
			request:     "/update/unknown/testUnknown/100",
			requestType: "POST",
			contentType: "text/plain",
			body:        nil,
			want: want{
				statusCode: http.StatusNotImplemented,
				body:       "incorrect metric type",
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
	}
	suite.runTableTests(testData)
}

//TestListHandler тестируем ValueHandler
func (suite *HandlersSuite) TestValueHandler() {
	testData := []tests{
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
				body:       "not found",
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
				body:       "not found",
			},
		},
	}
	suite.runTableTests(testData)
}

//TestListHandler тестируем UpdateHandler
func (suite *HandlersSuite) TestUpdateHandler() {
	testData := []tests{
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
			name:        "get correct gauge value",
			request:     "/value/unknown/testGauge",
			requestType: "GET",
			contentType: "text/plain",
			body:        nil,
			want: want{
				statusCode: http.StatusNotFound,
				body:       "not found",
			},
		},
	}
	suite.runTableTests(testData)
}

// TestUpdateValueJSONHandler тестируем UpdateJSONHandler ValueJSONHandler
func (suite *HandlersSuite) TestUpdateValueJSONHandler() {
	testData := []tests{
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
			name:        "set gauge with GET",
			request:     "/update/",
			requestType: "GET",
			contentType: "application/json",
			body:        []byte(`{"id":"jsonGauge","type":"gauge","value": 111}`),
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				body:       "",
			},
		}, {
			name:        "set gauge with incorrect value",
			request:     "/update/",
			requestType: "POST",
			contentType: "application/json",
			body:        []byte(`{"id":"jsonGauge","type":"gauge","value": "incorrect value"}`),
			want: want{
				statusCode: http.StatusInternalServerError,
				body:       "invalid deserialization",
			},
		},

		{
			name:        "set gauge with json - bad json",
			request:     "/update/",
			requestType: "POST",
			contentType: "application/json",
			body:        []byte(`{"id":"jsonGauge""type":"gauge""value": 111}`),
			want: want{
				statusCode: http.StatusInternalServerError,
				body:       "invalid deserialization",
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
		{
			name:        "value - bad json",
			request:     "/value/",
			requestType: "POST",
			contentType: "application/json",
			body:        []byte(`{"id":"jsonGauge""type":"gauge"}`),
			want: want{
				statusCode: http.StatusInternalServerError,
				body:       `invalid deserialization`,
			},
		},
		{
			name:        "test get counter with bad contentType",
			request:     "/value/",
			requestType: "POST",
			contentType: "text/plain",
			body:        []byte(`{"id":"jsonCounter","type":"counter"}`),
			want: want{
				statusCode: http.StatusInternalServerError,
				body:       `content type not support`,
			},
		},
		{
			name:        "update gauge with bad contentType",
			request:     "/update/",
			requestType: "POST",
			contentType: "text/plain",
			body:        []byte(`{"id":"jsonGauge","type":"gauge","value": 111}`),
			want: want{
				statusCode: http.StatusInternalServerError,
				body:       "content type not support",
			},
		},
		{
			name:        "get counter that doesn't exists",
			request:     "/value/",
			requestType: "POST",
			contentType: "application/json",
			body:        []byte(`{"id":"notFoundCounter","type":"counter"}`),
			want: want{
				statusCode: http.StatusNotFound,
				body:       `not found`,
			},
		},
	}
	suite.runTableTests(testData)
}

// TestUpdateValuesJSONHandler тестируем /updates/
func (suite *HandlersSuite) TestUpdateValuesJSONHandler() {
	testData := []tests{
		{
			name:        "set gauge and counter with json",
			request:     "/updates/",
			requestType: "POST",
			contentType: "application/json",
			body: []byte(`[{"id":"jsonGauge","type":"gauge","value": 111},
							      {"id":"jsonCounter","type":"counter","delta": 123}]`),
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
		{
			name:        "/updates/ with invalid json",
			request:     "/updates/",
			requestType: "POST",
			contentType: "application/json",
			body: []byte(`[{"id":"jsonGauge","type":"gauge","value": 111,},
						   {"id":"jsonCounter","type":"counter","delta": 123,}]`),
			want: want{
				statusCode: http.StatusInternalServerError,
				body:       "invalid deserialization",
			},
		},
		{
			name:        "set gauge and counter with json",
			request:     "/updates/",
			requestType: "POST",
			contentType: "text/plain",
			body: []byte(`[{"id":"jsonGauge","type":"gauge","value": 111},
							      {"id":"jsonCounter","type":"counter","delta": 123}]`),
			want: want{
				statusCode: http.StatusInternalServerError,
				body:       "content type not support",
			},
		},
	}
	suite.runTableTests(testData)
}

// TestUpdateValueJSONHandler тестируем UpdateJSONHandler ValueJSONHandler
func (suite *HandlersSuite) TestPingHandler() {
	testData := []tests{
		{
			name:        "check ping handler with fake storage",
			request:     "/ping",
			requestType: "GET",
			contentType: "",
			body:        nil,
			want: want{
				statusCode: http.StatusOK,
				body:       "ok",
			},
		},
	}
	suite.runTableTests(testData)
}

func newExampleServer() *httptest.Server {
	conf := config.Config{}
	repo := repository.NewRepository(conf.GeneralCfg())
	//! Это пустой вариант хранилища без состояние. Ошибок нет
	pStore, _ := storage.NewFakeStorage()

	r := NewRouter(repo, &conf, pStore)
	ts := httptest.NewServer(r)
	return ts
}

func ExampleHandler_Ping() {
	ts := newExampleServer()

	req, err := http.NewRequest("GET", ts.URL+"/ping", nil)
	if err != nil {
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	fmt.Printf("status = %d, body = %s\n", resp.StatusCode, string(respBody))
	// Output:
	// status = 200, body = ok
}

func ExampleHandler_UpdateListJSON() {
	ts := newExampleServer()

	req, _ := http.NewRequest("POST", ts.URL+"/updates/",
		bytes.NewBuffer([]byte(`[{"id":"jsonGauge","type":"gauge","value": 111},
								 {"id":"jsonCounter","type":"counter","delta": 123}]`)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	fmt.Printf("status = %d, body = %s\n", resp.StatusCode, string(respBody))
	// Output:
	// status = 200, body = ok
}

func ExampleHandler_UpdateJSON() {
	ts := newExampleServer()

	req, _ := http.NewRequest("POST", ts.URL+"/update/",
		bytes.NewBuffer([]byte(`{"id":"jsonGauge","type":"gauge","value": 111}`)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	fmt.Printf("status = %d, body = %s\n", resp.StatusCode, string(respBody))

	// Output:
	// status = 200, body = ok
}

func ExampleHandler_Update() {
	ts := newExampleServer()

	req, _ := http.NewRequest("POST", ts.URL+"/update/counter/testCounter/100", nil)
	req.Header.Set("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	fmt.Printf("status = %d, body = %s\n", resp.StatusCode, string(respBody))

	// Output:
	// status = 200, body = ok
}

func ExampleHandler_Value() {
	ts := newExampleServer()
	{
		// Добавляем testCounter
		req, _ := http.NewRequest("POST", ts.URL+"/update/counter/testCounter/100", nil)
		req.Header.Set("Content-Type", "text/plain")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return
		}
		resp.Body.Close()
	}

	req, err := http.NewRequest("GET", ts.URL+"/value/counter/testCounter", nil)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	fmt.Printf("status = %d, body = %s\n", resp.StatusCode, string(respBody))

	// Output:
	// status = 200, body = 100
}

func ExampleHandler_ValueJSON() {
	ts := newExampleServer()
	{
		// Добавляем testCounter
		req, _ := http.NewRequest("POST", ts.URL+"/update/counter/jsonCounter/100", nil)
		req.Header.Set("Content-Type", "text/plain")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return
		}
		resp.Body.Close()
	}

	req, _ := http.NewRequest("POST", ts.URL+"/value/",
		bytes.NewBuffer([]byte(`{"id":"jsonCounter","type":"counter","delta":123}`)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	fmt.Printf("status = %d, body = %s\n", resp.StatusCode, string(respBody))

	// Output:
	// status = 200, body = {"id":"jsonCounter","type":"counter","delta":100}
}
