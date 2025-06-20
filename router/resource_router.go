package router

import (
	"fmt"
	"github.com/lwm-galactic/mcp/context"
	"github.com/lwm-galactic/mcp/core/message"
	"github.com/lwm-galactic/mcp/core/resources"
	"github.com/lwm-galactic/mcp/handler"
	"net/http"
	"regexp"
	"strings"
)

// 路由项
type route struct {
	pattern    *regexp.Regexp
	resource   resources.Resource
	paramNames []string
}

// ResourceRegistry 资源注册器
type ResourceRegistry struct {
	routes []*route
}

func NewResourceRegistry() *ResourceRegistry {
	return &ResourceRegistry{
		routes: make([]*route, 0),
	}
}

// Register 注册一个资源
func (r *ResourceRegistry) Register(uriPattern string, res resources.Resource) {
	rePattern := regexp.MustCompile(`\{([^}]+)\}`)

	matches := rePattern.FindAllStringSubmatch(uriPattern, -1)
	var paramNames []string
	for _, m := range matches {
		paramNames = append(paramNames, m[1])
	}

	reStr := rePattern.ReplaceAllStringFunc(uriPattern, func(s string) string {
		name := strings.Trim(s, "{}")
		return fmt.Sprintf(`(?P<%s>[^/]+)`, name)
	})

	re := regexp.MustCompile(`^` + reStr + `$`)
	r.routes = append(r.routes, &route{
		pattern:    re,
		resource:   res,
		paramNames: paramNames,
	})
}

// Get 查找资源（按 URI）
func (r *ResourceRegistry) Get(uri string) (resources.Resource, bool) {
	res, _, ok := r.Match(uri)
	return res, ok
}

func (r *ResourceRegistry) Match(path string) (resources.Resource, map[string]string, bool) {
	for _, rt := range r.routes {
		matches := rt.pattern.FindStringSubmatch(path)
		if matches == nil {
			continue
		}

		params := make(map[string]string)
		for i, name := range rt.pattern.SubexpNames()[1:] {
			if i < len(matches)-1 {
				params[name] = matches[i+1]
			}
		}

		return rt.resource, params, true
	}
	return nil, nil, false
}

// makeReadResourceHandler 创建 read_resource 的处理函数
func makeReadResourceHandler(registry *ResourceRegistry) handler.RPCHandlerFunc {
	return func(params map[string]interface{}, w http.ResponseWriter, req message.Request) message.Response {
		resp := message.Response{}
		uri, ok := params["uri"].(string)
		if !ok || uri == "" {
			resp.Errorf(req.ID, message.NewError(message.InvalidRequest, "uri is required"))
			return resp
		}

		resource, exists := registry.Get(uri)
		if !exists {
			resp.Errorf(req.ID, message.NewError(message.ResourceNotFound, fmt.Sprintf("%s is not found", uri)))
			return resp
		}

		ctx := &context.ResourceContext{
			Request: req,
			Params:  params,
		}

		content, err := resource.GetContent(ctx)
		if err != nil {
			resp.Errorf(req.ID, message.NewError(message.InternalError, err.Error()))
			return resp
		}
		resp.Success(req.ID, content)
		return resp
	}
}

// makeListResourcesHandler 创建 list_resources 的处理函数
func makeListResourcesHandler(registry *ResourceRegistry) handler.RPCHandlerFunc {
	return func(_ map[string]interface{}, w http.ResponseWriter, req message.Request) message.Response {
		resources := make([]map[string]interface{}, 0, len(registry.routes))
		for _, route := range registry.routes {
			resources = append(resources, map[string]interface{}{
				"name":        route.resource.Name(),
				"description": route.resource.Description(),
				"type":        route.resource.Type(),
				"pattern":     route.pattern.String(),
			})
		}
		resp := message.Response{}
		resp.Success(req.ID, resources)
		return resp
	}
}

func RegisterResourceRoutes(mux *http.ServeMux, registry *ResourceRegistry) {
	prefix := "/resource/"

	mux.HandleFunc(prefix+"read", func(w http.ResponseWriter, r *http.Request) {
		handleRPC(w, r, makeReadResourceHandler(registry))
	})

	mux.HandleFunc(prefix+"list", func(w http.ResponseWriter, r *http.Request) {
		handleRPC(w, r, makeListResourcesHandler(registry))
	})
}
