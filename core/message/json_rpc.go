package message

// 符合 JSON RPC 格式的结构体 统一请求和返回
const (
	JsonRpcVersion = "2.0"
)

// Request 统一接收请求
type Request struct {
	ID      string                 `json:"id,omitempty"`      // 可选：用于追踪请求
	JSONRPC string                 `json:"jsonrpc"`           // 固定为 "2.0"
	Method  string                 `json:"method"`            // 方法名
	Params  map[string]interface{} `json:"params,omitempty"`  // 参数
	Session string                 `json:"session,omitempty"` // 会话标识
}

// Response 统一返回
type Response struct {
	JSONRPC string      `json:"jsonrpc"`          // 固定为 "2.0"
	ID      string      `json:"id"`               // 请求中的 id，用于匹配响应
	Result  interface{} `json:"result,omitempty"` // 正常结果
	Error   *Error      `json:"error,omitempty"`  // 错误信息
}

func Success(id string, data interface{}) *Response {
	return &Response{
		JSONRPC: JsonRpcVersion,
		ID:      id,
		Result:  data,
		Error:   nil,
	}
}

func Errorf(id string, code ErrorCode, message string, data ...string) *Response {
	return &Response{
		JSONRPC: JsonRpcVersion,
		ID:      id,
		Error:   newJSONRPCError(code, message, data...),
		Result:  nil,
	}
}
