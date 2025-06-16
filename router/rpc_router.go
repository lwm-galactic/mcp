package router

import (
	"encoding/json"
	"github.com/lwm-galactic/mcp/core/message"
	"github.com/lwm-galactic/mcp/handler"
	"net/http"
)

func handleRPC(w http.ResponseWriter, r *http.Request, handler handler.RPCHandlerFunc) {
	var req message.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "invalid_request", http.StatusBadRequest)
		return
	}

	if err := handler(req.Params, w, r); err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
	}
}
