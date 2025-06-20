package context

import (
	"github.com/lwm-galactic/mcp/core/message"
)

// ResourceContext 是传递给 ResourceHandler 的上下文
type ResourceContext struct {
	Request message.Request
	Params  map[string]interface{} // URL 参数
}
