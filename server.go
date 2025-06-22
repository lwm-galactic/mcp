package zeno

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"os"
)

type TransportType string

const (
	TransportSSE            TransportType = "sse"
	TransportStreamableHTTP TransportType = "streamable-http"
)

const (
	DefaultVersion = "v1.0.0"
)

type Server struct {
	Name    string
	Version string
	Prefix  string
	Log     *zap.Logger

	middlewares []func(http.Handler) http.Handler
}

func NewServer(name string) *Server {
	return &Server{
		Name:        name,
		Version:     DefaultVersion,
		middlewares: []func(http.Handler) http.Handler{},
	}
}

func (s *Server) Run(transport TransportType, addr ...string) error {

	if GetMode() == DebugMode {
		s.middlewares = append(s.middlewares, NewLoggingMiddleware(s.Log))
	}
	// 根据 transport 类型选择处理方式
	switch transport {
	case TransportSSE:
		return s.startSSE(resolveAddress(addr))
	case TransportStreamableHTTP:
		return s.startStreamableHTTP(resolveAddress(addr))
	default:
		return fmt.Errorf("unsupported transport type: %s", transport)
	}
}

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			debugPrint("Environment variable PORT=\"%s\"", port)
			return ":" + port
		}
		debugPrint("Environment variable PORT is undefined. Using port :8080 by default")
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too many parameters")
	}
}

func (s *Server) startSSE(addr string) error {
	return nil
}

func (s *Server) startStreamableHTTP(addr string) error {
	return nil
}
