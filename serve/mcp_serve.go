package serve

import (
	"fmt"
	"github.com/lwm-galactic/mcp/core/resources"
	"github.com/lwm-galactic/mcp/core/tools"
	"github.com/lwm-galactic/mcp/logger"
	"github.com/lwm-galactic/mcp/middleware"
	"github.com/lwm-galactic/mcp/router"
	"go.uber.org/zap"
	"net/http"
)

type TransportType string

const (
	TransportSSE            TransportType = "sse"
	TransportStreamableHTTP TransportType = "streamable-http"
)

type Mode string

const (
	DebugMode   Mode = "debug"
	ReleaseMode Mode = "release"
	TestMode    Mode = "test"
)

type McpServe struct {
	Name             string
	Describe         string
	Mode             Mode
	Log              *zap.Logger
	resourceRegistry *router.ResourceRegistry
	toolRegistry     *router.ToolRegistry
	middlewares      []func(http.Handler) http.Handler
	promptRegistry   *router.ToolRegistry
}

func NewMcpServe(name, describe string) *McpServe {
	return &McpServe{
		Name:             name,
		Describe:         describe,
		Mode:             DebugMode,
		Log:              logger.GetLogger(),
		toolRegistry:     router.NewToolRegistry(),
		resourceRegistry: router.NewResourceRegistry(),
		promptRegistry:   router.NewToolRegistry(),
		middlewares:      []func(http.Handler) http.Handler{},
	}
}

func (s *McpServe) SetMode(mode Mode) {
	s.Mode = mode
}

func (s *McpServe) Run(addr string, transport TransportType) error {
	if s.Mode == DebugMode {
		s.middlewares = append(s.middlewares, middleware.NewLoggingMiddleware(s.Log))
	}
	// 根据 transport 类型选择处理方式
	switch transport {
	case TransportSSE:
		return s.startSSE(addr)
	case TransportStreamableHTTP:
		return s.startStreamableHTTP(addr)
	default:
		return fmt.Errorf("unsupported transport type: %s", transport)
	}
}

// RegisterResource 注册资源
func (s *McpServe) RegisterResource(uriPattern string, resource resources.Resource) {
	s.resourceRegistry.Register(uriPattern, resource)
}

// RegisterTool 注册工具
func (s *McpServe) RegisterTool(schema tools.ToolSchema, tool tools.Tool) {
	s.toolRegistry.Register(tool, schema)
}

// Use 使用中间键
func (s *McpServe) Use(middleware func(http.Handler) http.Handler) {
	s.middlewares = append(s.middlewares, middleware)
}

func (s *McpServe) startSSE(addr string) error {
	// 创建多路复用器
	mux := http.NewServeMux()

	// 注册资源和工具路由
	router.RegisterResourceRoutes(mux, s.resourceRegistry)
	router.RegisterToolRoutesSSE(mux, s.toolRegistry)

	// 构建中间件链
	h := http.Handler(mux)
	for _, m := range s.middlewares {
		h = m(h)
	}
	return http.ListenAndServe(addr, h)
}

func (s *McpServe) startStreamableHTTP(addr string) error {
	// 创建多路复用器
	mux := http.NewServeMux()

	// 注册资源和工具路由
	router.RegisterResourceRoutes(mux, s.resourceRegistry)
	router.RegisterToolRoutesStreamableHTTP(mux, s.toolRegistry)

	// 构建中间件链
	h := http.Handler(mux)
	for _, m := range s.middlewares {
		h = m(h)
	}
	return http.ListenAndServe(addr, h)
}
