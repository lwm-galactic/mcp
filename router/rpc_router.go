package router

import (
	"encoding/json"
	"github.com/lwm-galactic/mcp/core/message"
	"github.com/lwm-galactic/mcp/handler"
	"net/http"
)

func handleRPC(w http.ResponseWriter, r *http.Request, handler handler.RPCHandlerFunc) {
	var req message.Request
	resp := message.Response{}
	flusher, ok := w.(http.Flusher)
	if !ok { // 不允许流式返回
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.Errorf(req.ID, message.NewError(message.InvalidRequest, "body is required"))
	} else {
		resp = handler(req.Params, w, req)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
	flusher.Flush()
}
