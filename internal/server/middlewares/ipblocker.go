package middlewares

import (
	"context"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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

func (b *IPBlocker) IsAllowIP(realIP string) bool {
	if b.cidr == nil {
		return true
	}
	clientIP := net.ParseIP(realIP)
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
		if b.IsAllowIP(r.Header.Get("X-Real-IP")) {
			next.ServeHTTP(w, r)
			return
		}
		w.WriteHeader(http.StatusForbidden)
	})
}

func IPBlock(cidr string) func(next http.Handler) http.Handler {
	blocker := NewIPBlocker(cidr)
	return blocker.Handler
}

func (b *IPBlocker) unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if b.cidr == nil {
		return handler(ctx, req)
	}
	var realIP string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("X-Real-IP")
		if len(values) > 0 {
			realIP = values[0]
			// Итого если в наличии есть параметр метадаты. И он не подходит, то PermissionDenied иначе - ок
			if !b.IsAllowIP(realIP) {
				return nil, status.Error(codes.PermissionDenied, "incorrect X-Real-IP")
			}
		}
	}
	return handler(ctx, req)
}

func IPBlockInterceptor(cidr string) grpc.UnaryServerInterceptor {
	blocker := NewIPBlocker(cidr)
	return blocker.unaryInterceptor
}
