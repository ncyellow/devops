package middlewares

import (
	"net"
	"net/http"
)

type IPBlocker struct {
	cidr *net.IPNet
}

func NewIPBlocker(cidr string) *IPBlocker {
	if cidr == "" {
		return &IPBlocker{
			cidr: nil,
		}
	}
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return &IPBlocker{
			cidr: nil,
		}
	}
	return &IPBlocker{
		cidr: ipNet,
	}
}

func (b *IPBlocker) IsAllowIP(realIp string) bool {
	if b.cidr == nil {
		return true
	}
	clientIP := net.ParseIP(realIp)
	return b.cidr.Contains(clientIP)
}

// Handler middleware для фильтрации IP адресов
func (b *IPBlocker) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if b.cidr == nil {
			next.ServeHTTP(w, r)
			return
		}

		// проверяем, что у клиента передан заголовок gzip-сжатие
		if !b.IsAllowIP(r.Header.Get("X-Real-IP")) {
			next.ServeHTTP(w, r)
			return
		}
		w.WriteHeader(http.StatusForbidden)
		return
	})
}

func IPBlock(cidr string) func(next http.Handler) http.Handler {
	blocker := NewIPBlocker(cidr)
	return blocker.Handler
}
