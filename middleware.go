package zeno

import (
	"bytes"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

// LoggingMiddleware 使用 zap 记录请求日志
type LoggingMiddleware struct {
	logger *zap.Logger
}

func NewLoggingMiddleware(logger *zap.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// 读取请求体（注意：body 是一次性读取的，后续处理需要重新设置 Body）
			body, err := io.ReadAll(r.Body)
			if err != nil {
				logger.Error("Failed to read request body", zap.Error(err))
			}
			// 重新设置 Body，以便后续处理器能正常读取
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			// 构建日志字段
			logFields := []zap.Field{
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("query", r.URL.RawQuery),
				zap.String("user_agent", r.UserAgent()),
				zap.String("content_type", r.Header.Get("Content-Type")),
				zap.ByteString("body", body),
			}

			logger.Info("Request received", logFields...)

			// 调用下一个处理器
			next.ServeHTTP(w, r)

			// 可选：记录响应结束时间或状态码
			// latency := time.Since(start)
			// logger.Debug("Request completed", zap.Duration("latency", latency))
		})
	}
}
