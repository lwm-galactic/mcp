package router

import (
	"fmt"
	"github.com/lwm-galactic/mcp/core/message"
	"github.com/lwm-galactic/mcp/core/tools"
	"github.com/lwm-galactic/mcp/handler"
	"net/http"
)

// MCPTOOL 定义一个符合 MCP 协议的工具描述结构体
type MCPTOOL struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
	Annotations map[string]interface{} `json:"annotations,omitempty"`
}

type ToolRegistry struct {
	tools   map[string]tools.Tool
	schemas map[string]tools.ToolSchema
}

func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools:   make(map[string]tools.Tool),
		schemas: make(map[string]tools.ToolSchema),
	}
}

// Register 注册一个工具及其元数据
func (r *ToolRegistry) Register(tool tools.Tool, schema tools.ToolSchema) {
	r.tools[tool.Name()] = tool
	r.schemas[tool.Name()] = schema
}

// Get 获取工具
func (r *ToolRegistry) Get(name string) (tools.Tool, bool) {
	tool, ok := r.tools[name]
	return tool, ok
}

// GetSchema 获取工具元数据
func (r *ToolRegistry) GetSchema(name string) (tools.ToolSchema, bool) {
	schema, ok := r.schemas[name]
	return schema, ok
}

// ListSchemas 列出所有工具元数据
func (r *ToolRegistry) ListSchemas() []tools.ToolSchema {
	schemas := make([]tools.ToolSchema, 0, len(r.schemas))
	for _, s := range r.schemas {
		schemas = append(schemas, s)
	}
	return schemas
}

func makeStreamableInvokeToolHandler(registry *ToolRegistry) handler.RPCHandlerFunc {
	return func(params map[string]interface{}, w http.ResponseWriter, req message.Request) message.Response {
		var resp message.Response
		toolName, ok := params["tool_name"].(string)
		if !ok || toolName == "" { // 工具名称未传
			resp.Errorf(req.ID, message.NewError(message.InvalidRequest, "tool_name is required"))
			return resp
		}

		args, ok := params["arguments"].(map[string]interface{})
		if !ok { // 工具参数错误
			resp.Errorf(req.ID, message.NewError(message.ParseError, "params arguments is required"))
			return resp
		}

		tool, exists := registry.Get(toolName)
		if !exists { // 工具未找到
			resp.Errorf(req.ID, message.NewError(message.MethodNotFound, fmt.Sprintf("%s method is not found", toolName)))
			return resp
		}

		result, err := tool.Execute(args)
		if err != nil { // 工具执行错误
			resp.Errorf(req.ID, message.NewError(message.InternalError, fmt.Sprintf("method execute error %s", err.Error())))
			return resp
		}
		resp.Success(req.ID, result)
		return resp
	}
}

func isStreamingRequest(r *http.Request) bool {
	return r.Header.Get("Accept") == "text/event-stream"
}

// makeListToolsHandler 创建 listTools 的处理函数
func makeListToolsHandler(registry *ToolRegistry) handler.RPCHandlerFunc {
	return func(_ map[string]interface{}, w http.ResponseWriter, req message.Request) message.Response {
		schemas := registry.ListSchemas()

		var mcpTools []MCPTOOL
		for _, schema := range schemas {
			params := make([]map[string]interface{}, len(schema.Parameters))
			for i, p := range schema.Parameters {
				params[i] = map[string]interface{}{
					"name":        p.Name,
					"type":        p.Type,
					"description": p.Description,
					"required":    p.Required,
				}
			}

			mcpTools = append(mcpTools, MCPTOOL{
				Name:        schema.Metadata.Name,
				Description: schema.Metadata.Description,
				InputSchema: map[string]interface{}{
					"properties": params,
				},
				Annotations: schema.Annotations,
			})
		}

		resp := message.Response{}
		resp.Success(req.ID, mcpTools)
		return resp
	}
}

func RegisterToolRoutesSSE(mux *http.ServeMux, registry *ToolRegistry) {

	mux.HandleFunc("/tool/invoke", func(w http.ResponseWriter, r *http.Request) {
		// 设置 SSE 响应头
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		handleRPC(w, r, makeStreamableInvokeToolHandler(registry))
	})

	mux.HandleFunc("/tool/list", func(w http.ResponseWriter, r *http.Request) {
		// 设置 SSE 响应头
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		handleRPC(w, r, makeListToolsHandler(registry))
	})

	mux.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {

	})
}

func RegisterToolRoutesStreamableHTTP(mux *http.ServeMux, registry *ToolRegistry) {
	mux.HandleFunc("/tool/invoke", func(w http.ResponseWriter, r *http.Request) {
		// 设置 headers
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		handleRPC(w, r, makeStreamableInvokeToolHandler(registry))
	})

	mux.HandleFunc("/tool/list", func(w http.ResponseWriter, r *http.Request) {
		// 设置 headers
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		handleRPC(w, r, makeStreamableInvokeToolHandler(registry))
	})
}
