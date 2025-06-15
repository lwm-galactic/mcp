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
	// 创建多路复用器
	mux := http.NewServeMux()

	// 注册资源和工具路由
	router.RegisterResourceRoutes(mux, s.resourceRegistry, s.resourcePrefix)
	router.RegisterToolRoutes(mux, s.toolRegistry)

	// 构建中间件链
	h := http.Handler(mux)
	for _, middleware := range s.middlewares {
		h = middleware(h)
	}

	// 根据 transport 类型选择处理方式
	switch transport {
	case TransportSSE, TransportHTTP, TransportStreamableHTTP:
		// 暂时统一使用标准 HTTP 处理
		return http.ListenAndServe(addr, h)
	default:
		return fmt.Errorf("unsupported transport type: %s", transport)
	}
}

// RegisterResource 注册资源
func (s *McpServe) RegisterResource(uriPattern string, handler handler.ResourceHandler) {
	s.resourceRegistry.RegisterResource(uriPattern, handler)
}

func (s *McpServe) RegisterTool(schema tools.ToolSchema, tool tools.Tool) {
	s.toolRegistry.Register(tool, schema)
}

// RegisterTool 注册工具

func (s *McpServe) Use(middleware func(http.Handler) http.Handler) {
	s.middlewares = append(s.middlewares, middleware)
}
