package zeno

import (
	"fmt"
	"github.com/lwm-galactic/logger"
	"net/http"
)

// client 是每个 SSE 客户端的抽象
type client struct {
	id   string
	send chan string // 发送消息的通道
}

func handleSSE(server *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientID := r.URL.Query().Get("id")
		if clientID == "" {
			http.Error(w, "missing client id", http.StatusBadRequest)
			return
		}

		server.register(clientID, w, r)
	}
}

// writePump 处理消息写入客户端
func (s *Server) writePump(client *client, w http.ResponseWriter, r *http.Request) {
	defer func() {
		s.unregister(client.id)
	}()

	for {
		select {
		case message, ok := <-client.send:
			if !ok {
				// 通道被关闭
				return
			}
			_, err := fmt.Fprintf(w, "data: %s\n\n", message)
			if err != nil {
				logger.Errorf("Write error: %v", err)
				return
			}
			flusher, ok := w.(http.Flusher)
			if !ok {
				logger.Error("Streaming unsupported!")
				return
			}
			flusher.Flush()
		case <-r.Context().Done():
			// 客户端断开连接
			return
		}
	}
}

// Unregister 注销一个客户端
func (s *Server) unregister(clientID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	client, exists := s.clients[clientID]
	if !exists {
		return
	}

	close(client.send)
	delete(s.clients, clientID)
	logger.Infof("Client %s unregistered", clientID)
}

// Register 注册一个新的 SSE 客户端
func (s *Server) register(clientID string, w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.clients[clientID]; exists {
		logger.Errorf("Client %s already exists", clientID)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 创建客户端
	client := &client{
		id:   clientID,
		send: make(chan string, 100), // 带缓冲的通道
	}

	s.clients[clientID] = client

	// 启动写入协程
	go s.writePump(client, w, r)
}
