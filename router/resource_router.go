package router

import (
	"fmt"
	"github.com/lwm-galactic/mcp/context"
	"github.com/lwm-galactic/mcp/handler"
	"net/http"
	"regexp"
	"strings"
)

// ResourceRegistry 资源注册器
type ResourceRegistry struct {
	routes []*route
}

func NewResourceRegistry() *ResourceRegistry {
	return &ResourceRegistry{
		routes: make([]*route, 0),
	}
}

// RegisterResource 注册一个资源处理器
func (r *ResourceRegistry) RegisterResource(uriPattern string, handler handler.ResourceHandler) {
	// 替换参数部分为正则表达式组
	rePattern := regexp.MustCompile(`\{([^}]+)\}`)

	matches := rePattern.FindAllStringSubmatch(uriPattern, -1)
	var paramNames []string
	for _, m := range matches {
		paramNames = append(paramNames, m[1])
	}

	// 将 {name} 替换为命名捕获组 (?P<name>[^/]+)
	reStr := rePattern.ReplaceAllStringFunc(uriPattern, func(s string) string {
		name := strings.Trim(s, "{}")
		return fmt.Sprintf(`(?P<%s>[^/]+)`, name)
	})

	// 构建正则表达式
	re := regexp.MustCompile(`^` + reStr + `$`)
	r.routes = append(r.routes, &route{
		pattern:    re,
		handler:    handler,
		paramNames: paramNames,
	})
}

// Match 查找匹配的路由和参数
func (r *ResourceRegistry) Match(path string) (handler.ResourceHandler, map[string]string, bool) {
	for _, route := range r.routes {
		matches := route.pattern.FindStringSubmatch(path)
		if matches == nil {
			continue
		}

		params := make(map[string]string)
		for i, name := range route.pattern.SubexpNames()[1:] {
			if i < len(matches)-1 {
				params[name] = matches[i+1]
			}
		}

		return route.handler, params, true
	}
	return nil, nil, false
}

// RegisterResourceRoutes 注册资源访问路由
func RegisterResourceRoutes(mux *http.ServeMux, registry *ResourceRegistry, prefix string) {
	if prefix == "" {
		prefix = "/resource/"
	}

	mux.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
		// 提取路径部分
		path := r.URL.Path[len(prefix):]
		if path == "" {
			sendError(w, "missing resource path", http.StatusBadRequest)
			return
		}

		// 匹配路由并获取处理器和参数
		handlerFunc, params, ok := registry.Match(path)
		if !ok {
			sendError(w, fmt.Sprintf("resource '%s' not found", path), http.StatusNotFound)
			return
		}

		// 构造上下文
		ctx := &context.ResourceContext{
			Request: r,
			Params:  params,
		}

		// 执行资源处理器
		data, err := handlerFunc(ctx)
		if err != nil {
			sendError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 返回结果
		sendSuccess(w, map[string]interface{}{
			"path":   path,
			"params": params,
			"data":   string(data),
		})
	})
}

func RegisterResourceListRoute(mux *http.ServeMux, registry *ResourceRegistry, prefix string) {
	mux.HandleFunc(prefix+"list", func(w http.ResponseWriter, r *http.Request) {
		var routesInfo []map[string]interface{}
		for _, route := range registry.routes {
			routesInfo = append(routesInfo, map[string]interface{}{
				"pattern":    route.pattern.String(),
				"paramNames": route.paramNames,
			})
		}
		sendSuccess(w, map[string]interface{}{
			"resources": routesInfo,
		})
	})
}
