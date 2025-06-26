package zeno

import (
	"fmt"
	"github.com/lwm-galactic/logger"
	"github.com/lwm-galactic/zeno/core/resources"
	"github.com/lwm-galactic/zeno/core/tools"
	"sync"

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
	name        string
	version     string
	prefix      string
	mode        string
	router      *rpcRouter
	clients     map[string]*client // 所有连接的客户端
	mu          sync.Mutex         // sse 链接池互斥锁，保证并发安全
	middlewares []func(http.Handler) http.Handler
}

func NewServer(name string) *Server {
	return &Server{
		name:        name,
		version:     DefaultVersion,
		router:      newRPCRouter(),
		middlewares: []func(http.Handler) http.Handler{},
	}
}

func (s *Server) Run(transport TransportType, addr ...string) error {
	if GetMode() == DebugMode {
		s.middlewares = append(s.middlewares, NewRequestLoggingMiddleware())
		s.printTool()
		s.printTool()
	}
	if GetMode() == ReleaseMode {
		logger.SetLogLevel(-2)
	}
	// 根据 transport 类型选择处理方式
	switch transport {
	case TransportSSE:
		return s.startSSE(resolveAddress(addr))
	case TransportStreamableHTTP:
		return s.startStreamableHTTP(resolveAddress(addr))
	default:
		logger.Errorf("Unsupported transport type: %s", transport)
		return fmt.Errorf("unsupported transport type: %s", transport)
	}
}

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			logger.Debug("Environment variable PORT=\"%s\"", port)
			return ":" + port
		}
		logger.Debug("Environment variable PORT is undefined. Using port :8080 by default")
		return ":8080"
	case 1:
		return addr[0]
	default:
		logger.Panic("too many parameters")
		return ""
	}
}

func (s *Server) startSSE(addr string) error {
	s.clients = make(map[string]*client) // 初始化链接池
	http.HandleFunc("/sse", handleSSE(s))
	return nil
}

func (s *Server) startStreamableHTTP(addr string) error {
	return nil
}

func (s *Server) RegisterTool(tool tools.Tool) {
	s.router.registerTool(tool)
}

func (s *Server) RegisterResource(resource resources.Resource) {
	s.router.registerResource(resource)
}

func init() {
	logger.SetModName("[zeno]")
}

// 打印 注册的工具列表
func (s *Server) printTool() {
	info := "Tool: \n"
	for _, tool := range s.router.toolList {
		info += fmt.Sprintf("[zeno]\t Name:%s \t --> %s \n ", tool.Name(), tool.Description())
	}
	logger.Debug(info)
}

// 打印 注册的资源列表
func (s *Server) printResource() {
	info := "Resource: \n"
	for _, resource := range s.router.resourceList {
		info += fmt.Sprintf("[zeno]\t Name:%s \t --> %s \n ", resource.Name(), resource.Description())
	}
	logger.Debug(info)
}
