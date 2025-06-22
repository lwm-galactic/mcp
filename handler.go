package zeno

import (
	"encoding/json"
	"fmt"
	"github.com/lwm-galactic/zeno/core/message"
	"net/http"
)

type HandlerFunc func(*Context)

type RPCHandlerFunc func(params map[string]interface{}, w http.ResponseWriter, r message.Request) message.Response

func RPCHandler(w http.ResponseWriter, ctx *Context, handler RPCHandlerFunc) {
	var req message.Request
	fmt.Println(req)
	resp := message.Response{}
	flusher, ok := w.(http.Flusher)
	if !ok { // 不允许流式返回
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
	flusher.Flush()
}
