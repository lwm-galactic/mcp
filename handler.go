package zeno

import (
	"github.com/lwm-galactic/zeno/core/message"
)

type HandlerFunc func(*Context)

type rpcHandlerFunc func(message.Request) message.Response

/*
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


*/
