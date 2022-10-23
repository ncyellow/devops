package middlewares

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// getLocalCIDR возвращает CIDR текущей сети
func getLocalCIDR() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.String()
			}
		}
	}
	return ""
}

// getLocalCIDRWithCustomMask возвращает CIDR текущей сети
func getLocalCIDRWithCustomMask(customMask string) string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return fmt.Sprintf("%s/%s", ipnet.IP.String(), customMask)
			}
		}
	}
	return ""
}

func TestIPBlock(t *testing.T) {

	tests := map[string]struct {
		request          func() *http.Request
		cidr             string
		expectedResponse string
		expectedCode     int
	}{
		"IP has benn blocked": {
			func() *http.Request {
				req, _ := http.NewRequest("GET", "/", nil)
				req.Header.Add("X-Real-IP", getLocalCIDR())
				return req
			},
			getLocalCIDR(),
			"IP was been blocked",
			http.StatusForbidden,
		},
		"IP hasn't been blocked because we dont have X-REAL-IP": {
			func() *http.Request {
				req, _ := http.NewRequest("GET", "/", nil)
				return req
			},
			getLocalCIDR(),
			"ok",
			http.StatusOK,
		},

		"IP is ok": {
			func() *http.Request {
				req, _ := http.NewRequest("GET", "/", nil)
				req.Header.Add("X-Real-IP", getLocalCIDR())
				return req
			},
			getLocalCIDRWithCustomMask("255"),
			"ok",
			http.StatusOK,
		},
		"cidr is empty": {
			func() *http.Request {
				req, _ := http.NewRequest("GET", "/", nil)
				req.Header.Add("X-Real-IP", getLocalCIDR())
				return req
			},
			"",
			"ok",
			http.StatusOK,
		},
		"cidr is not valid": {
			func() *http.Request {
				req, _ := http.NewRequest("GET", "/", nil)
				req.Header.Add("X-Real-IP", getLocalCIDR())
				return req
			},
			"12343",
			"ok",
			http.StatusOK,
		},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()

		r := chi.NewRouter()
		r.Use(IPBlock(test.cidr))

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			response := "ok"
			w.Write([]byte(response))
		})
		r.ServeHTTP(w, test.request())

		assert.Equal(t, w.Code, test.expectedCode)
		assert.Equal(t, w.Body.String(), test.expectedResponse)
	}
}
