package resources

import "github.com/lwm-galactic/mcp/context"

type ResourceType string

const (
	ResourceTypeText ResourceType = "text"
	ResourceTypeJSON ResourceType = "json"
	ResourceTypeFile ResourceType = "file"
	ResourceTypeURL  ResourceType = "url"
)

// Resource resources.go - 接口定义
type Resource interface {
	Name() string
	Description() string
	Type() ResourceType // 可以是 "text", "json", "file", "url" 等
	GetContent(ctx *context.ResourceContext) (interface{}, error)
}
