package zeno

import (
	"github.com/lwm-galactic/logger"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

// RequestLoggingMiddleware 记录请求日志 中间件
type RequestLoggingMiddleware struct {
}

// NewRequestLoggingMiddleware 创建一个新的请求日志中间件
func NewRequestLoggingMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 记录开始时间
			start := time.Now()

			// 包装 ResponseWriter 以捕获状态码
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// 调用下一个 handler
			next.ServeHTTP(rw, r)

			// 计算耗时
			duration := time.Since(start)

			// 打印日志
			logger.Info("%s %s %s %d %v",
				r.RemoteAddr,
				r.Method,
				r.URL.Path,
				rw.statusCode,
				duration,
			)
		})
	}
}

// responseWriter 是一个包装器，用于捕获写入的状态码
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader 捕获状态码
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write 实现 Write 方法
func (rw *responseWriter) Write(body []byte) (int, error) {
	return rw.ResponseWriter.Write(body)
}

// Header 实现 Header 方法
func (rw *responseWriter) Header() http.Header {
	return rw.ResponseWriter.Header()
}
