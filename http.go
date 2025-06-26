package zeno

import (
	"encoding/json"
	"fmt"
	"github.com/lwm-galactic/logger"
	"github.com/lwm-galactic/zeno/core/message"
	"github.com/lwm-galactic/zeno/core/resources"
	"github.com/lwm-galactic/zeno/core/tools"
	"io"
	"net/http"
	"strings"
)

// sse 和 流式http 启动都需要的启动，用于接收符合http-rpc协议的 post 请求
func (s *Server) serverStart(addr string) error {
	mux := http.NewServeMux()

	mux.Handle(fmt.Sprintf("/%s/", s.prefix), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if requiresAcceptTypes(r, "application/json", "text/event-stream") {
			// 无效请求

		}

		if requireContentType(r, "application/json") {
			// 无效请求
		}

		// 读取请求体
		body, err := io.ReadAll(r.Body)
		if err != nil {
			// body 违规
		}
		defer r.Body.Close() // 记得关闭 Body
		// 解析 JSON 到结构体
		var req message.Request
		if err := json.Unmarshal(body, &req); err != nil {
			// json 序列化失败
		}

	}))

	logger.Info("Server is starting at http://%s ", addr)
	// 启动 HTTP 服务
	return http.ListenAndServe(addr, mux)
}

func requiresAcceptTypes(r *http.Request, required ...string) bool {
	// 检查指定 Header // Accept 必须 有 application/json, text/event-stream
	accept := r.Header.Get("Accept")
	if accept == "" {
		return false
	}

	requiredSet := make(map[string]bool)
	for _, t := range required {
		requiredSet[t] = false
	}

	types := strings.Split(accept, ",")
	for _, t := range types {
		mediaType := strings.TrimSpace(strings.Split(t, ";")[0])
		if _, exists := requiredSet[mediaType]; exists {
			requiredSet[mediaType] = true
		}
	}

	for _, seen := range requiredSet {
		if !seen {
			return false
		}
	}
	return true
}

func requireContentType(r *http.Request, expectedType string) bool {
	contentType := r.Header.Get("Content-Type")
	if contentType != expectedType {
		return false
	}
	return true
}

type rpcRouter struct {
	handlers     map[string]rpcHandlerFunc
	toolMap      map[string]tools.Tool
	resourceMap  map[string]resources.Resource
	toolList     []tools.Tool
	resourceList []resources.Resource
}

func newRPCRouter() *rpcRouter {
	r := &rpcRouter{
		handlers: make(map[string]rpcHandlerFunc),
	}

	// 初始化固定路由
	r.handlers[Initialize] = r.initHandler
	r.handlers[ToolsList] = r.toolsListHandler
	r.handlers[ToolsCall] = r.toolsCallHandler

	return r
}

const (
	Initialize = "initialize" //初始化请求 客户端发送具有协议版本和功能的 initialize 请求
	ToolsList  = "tools/list" //请求获取工具列表
	ToolsCall  = "tools/call" //请求获取工具执行
)

func (r *rpcRouter) rpcHandler(request message.Request) message.Response {
	if handler, exists := r.handlers[request.Method]; exists {
		return handler(request)
	}
	// 默认返回错误响应
	return message.Response{
		Error: &message.Error{Message: "Method not found"},
	}
}

func (r *rpcRouter) registerTool(tool tools.Tool) {
	if _, ok := r.toolMap[tool.Name()]; ok {
		logger.Errorf("Tool %s already registered", tool.Name())
	}
	r.toolMap[tool.Name()] = tool
	r.toolList = append(r.toolList, tool)
}

func (r *rpcRouter) registerResource(resource resources.Resource) {
	if _, ok := r.resourceMap[resource.Name()]; ok {
		logger.Errorf("Resource %s already registered", resource.Name())
	}
	r.resourceMap[resource.Name()] = resource
	r.resourceList = append(r.resourceList, resource)
}
func (r *rpcRouter) initHandler(request message.Request) message.Response {
	return message.Response{Result: "Initialized"}
}

func (r *rpcRouter) toolsListHandler(request message.Request) message.Response {
	// 获取工具列表
	return message.Response{Result: []string{"tool1", "tool2"}}
}

func (r *rpcRouter) toolsCallHandler(request message.Request) message.Response {
	// r.toolMap[request.Params["name"].string]
	return message.Response{Result: "Tool executed"}
}
