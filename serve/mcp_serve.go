package serve

import (
	"fmt"
	"github.com/lwm-galactic/mcp/core/tools"
	"github.com/lwm-galactic/mcp/handler"
	"github.com/lwm-galactic/mcp/router"
	"net/http"
	"strings"
)

type TransportType string

const (
	TransportSSE            TransportType = "sse"
	TransportHTTP           TransportType = "http"
	TransportStreamableHTTP TransportType = "streamable-http"
)

type McpServe struct {
	Name             string
	Describe         string
	resourceRegistry *router.ResourceRegistry
	toolRegistry     *router.ToolRegistry
	middlewares      []func(http.Handler) http.Handler
	resourcePrefix   string // 新增字段
}

func NewMcpServe(name, describe string) *McpServe {
	return &McpServe{
		Name:             name,
		Describe:         describe,
		toolRegistry:     router.NewToolRegistry(),
		resourceRegistry: router.NewResourceRegistry(),
		middlewares:      []func(http.Handler) http.Handler{},
		resourcePrefix:   "/resource/",
	}
}

// SetResourcePrefix 设置资源访问的 URL 前缀
func (s *McpServe) SetResourcePrefix(prefix string) {
	if prefix == "" {
		prefix = "/resource/"
	}
	// 确保以 '/' 结尾
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}
	s.resourcePrefix = prefix
}
func (s *McpServe) Run(addr string, transport TransportType) error {
	// 根据 transport 类型选择处理方式
	switch transport {
	case TransportSSE:
		return s.startSSE(addr)
	case TransportStreamableHTTP:
		return s.startStreamableHTTP(addr)
	case TransportHTTP:
		return s.startHTTP(addr)

	default:
		return fmt.Errorf("unsupported transport type: %s", transport)
	}
}

// RegisterResource 注册资源
func (s *McpServe) RegisterResource(uriPattern string, handler handler.ResourceHandler) {
	s.resourceRegistry.RegisterResource(uriPattern, handler)
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
	router.RegisterResourceRoutes(mux, s.resourceRegistry, s.resourcePrefix)
	router.RegisterToolRoutesSSE(mux, s.toolRegistry)

	// 构建中间件链
	h := http.Handler(mux)
	for _, middleware := range s.middlewares {
		h = middleware(h)
	}
	return http.ListenAndServe(addr, h)
}

func (s *McpServe) startHTTP(addr string) error {
	// 创建多路复用器
	mux := http.NewServeMux()

	// 注册资源和工具路由
	router.RegisterResourceRoutes(mux, s.resourceRegistry, s.resourcePrefix)
	router.RegisterToolRoutesHTTP(mux, s.toolRegistry)

	// 构建中间件链
	h := http.Handler(mux)
	for _, middleware := range s.middlewares {
		h = middleware(h)
	}
	return http.ListenAndServe(addr, h)
}

func (s *McpServe) startStreamableHTTP(addr string) error {
	// 创建多路复用器
	mux := http.NewServeMux()

	// 注册资源和工具路由
	router.RegisterResourceRoutes(mux, s.resourceRegistry, s.resourcePrefix)
	router.RegisterToolRoutesStreamableHTTP(mux, s.toolRegistry)

	// 构建中间件链
	h := http.Handler(mux)
	for _, middleware := range s.middlewares {
		h = middleware(h)
	}
	return http.ListenAndServe(addr, h)
}
