package router

import (
	"encoding/json"
	"github.com/lwm-galactic/mcp/core/message"
	"github.com/lwm-galactic/mcp/handler"
	"net/http"
	"regexp"
)

const (
	MethodGetResource   = "getResource"
	MethodInvokeTool    = "invokeTool"
	MethodListResources = "listResources"
	MethodListTools     = "listTools"
)

// 路由项
type route struct {
	pattern    *regexp.Regexp
	handler    handler.ResourceHandler
	paramNames []string
}

func sendSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(&message.Response{
		Status: "success",
		Data:   data,
	})
}

func sendError(w http.ResponseWriter, errorMsg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(&message.Response{
		Status: "error",
		Error:  errorMsg,
	})
}
