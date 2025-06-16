package handler

import "net/http"

// RPCHandler 定义统一的 RPC 处理器函数
type RPCHandler func(params map[string]interface{}) (interface{}, error)
type RPCHandlerFunc func(params map[string]interface{}, w http.ResponseWriter, r *http.Request) error
