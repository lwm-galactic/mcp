package context

import "net/http"

// ResourceContext 是传递给 ResourceHandler 的上下文
type ResourceContext struct {
	Request *http.Request
	Params  map[string]string // URL 参数
}
