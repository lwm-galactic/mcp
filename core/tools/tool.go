package tools

// Tool 工具接口
type Tool interface {
	Name() string
	Description() string
	Parameters() []ParamSchema
	Execute(params map[string]interface{}) (interface{}, error)
}

// ToolMetadata 工具元数据 结构体用于注册、序列化、展示
type ToolMetadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ParamSchema 描述单个参数的元信息
type ParamSchema struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"` // 如 "string", "integer", "boolean"
	Description string      `json:"description"`
	Default     interface{} `json:"default,omitempty"`
	Required    bool        `json:"required"`
}

// ToolSchema 描述整个工具的元信息
type ToolSchema struct {
	Metadata   ToolMetadata  `json:"metadata"`
	Parameters []ParamSchema `json:"parameters"`
}
