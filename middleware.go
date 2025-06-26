package zeno

import (
	"github.com/lwm-galactic/logger"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

// RequestLoggingMiddleware 记录请求日志 中间件
type RequestLoggingMiddleware struct {
}

func NewRequestLoggingMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		})
	}
}
