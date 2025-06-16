package router

import (
	"encoding/json"
	"github.com/lwm-galactic/mcp/core/message"
	"net/http"
)

const (
	MethodGetResource   = "getResource"
	MethodInvokeTool    = "invokeTool"
	MethodListResources = "listResources"
	MethodListTools     = "listTools"
)

func sendError(w http.ResponseWriter, msg string, code int) {
	resp := message.Response{
		Status: "error",
		Error:  msg,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(resp)
}

func sendSuccess(w http.ResponseWriter, data map[string]interface{}) {
	resp := message.Response{
		Status: "success",
		Data:   data,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
