package router

import (
	"encoding/json"
	"github.com/lwm-galactic/mcp/core/message"
	"github.com/lwm-galactic/mcp/core/tools"
	"net/http"
)

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

func RegisterToolRoutes(mux *http.ServeMux, registry *ToolRegistry) {
	mux.HandleFunc("/tool/invoke", func(w http.ResponseWriter, r *http.Request) {
		var req message.Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, "invalid_request", http.StatusBadRequest)
			return
		}

		if req.Method != MethodInvokeTool {
			sendError(w, "unknown_method", http.StatusBadRequest)
			return
		}

		toolName, ok := req.Params["tool_name"].(string)
		if !ok || toolName == "" {
			sendError(w, "missing_tool_name", http.StatusBadRequest)
			return
		}

		args, ok := req.Params["arguments"].(map[string]interface{})
		if !ok {
			sendError(w, "invalid_arguments", http.StatusBadRequest)
			return
		}

		tool, exists := registry.Get(toolName)
		if !exists {
			sendError(w, "tool_not_found", http.StatusNotFound)
			return
		}

		result, err := tool.Execute(args)
		if err != nil {
			sendError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sendSuccess(w, map[string]interface{}{
			"result": result,
		})
	})

	mux.HandleFunc("/tool/list", func(w http.ResponseWriter, r *http.Request) {
		schemas := registry.ListSchemas()
		sendSuccess(w, map[string]interface{}{
			"tools": schemas,
		})
	})
}
