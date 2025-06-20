package message

import (
	"encoding/json"
	"fmt"
)

type ErrorCode int

const (
	ParseError           ErrorCode = -32700 // 解析错误 [[1]]
	InvalidRequest       ErrorCode = -32600 // 无效请求 [[1]]
	MethodNotFound       ErrorCode = -32601 // 方法未找到 [[1]]
	InvalidParams        ErrorCode = -32602 // 参数无效 [[1]]
	InternalError        ErrorCode = -32603 // 内部错误 [[1]]
	ServerNotInitialized ErrorCode = -32001 // 服务器未初始化（自定义）
	ResourceNotFound     ErrorCode = -32002 // 资源未找到（自定义）
	UnsupportedMethod    ErrorCode = -32003 // 不支持的方法（自定义）
)

// Request 统一接收请求
type Request struct {
	ID      string                 `json:"id,omitempty"`      // 可选：用于追踪请求
	JSONRPC string                 `json:"jsonrpc"`           // 固定为 "2.0"
	Method  string                 `json:"method"`            // 方法名
	Params  map[string]interface{} `json:"params,omitempty"`  // 参数
	Session string                 `json:"session,omitempty"` // 会话标识
}

func (r *Request) UnmarshalParams(target interface{}) error {
	if r.Params == nil {
		return fmt.Errorf("empty params")
	}
	data, _ := json.Marshal(r.Params)
	return json.Unmarshal(data, target)
}

// Response 统一返回
type Response struct {
	JSONRPC string      `json:"jsonrpc"`          // 固定为 "2.0"
	ID      string      `json:"id"`               // 请求中的 id，用于匹配响应
	Result  interface{} `json:"result,omitempty"` // 正常结果
	Error   *Error      `json:"error,omitempty"`  // 错误信息
}
type Error struct {
	Code    ErrorCode   `json:"code"`           // 错误码
	Message string      `json:"message"`        // 错误描述
	Data    interface{} `json:"data,omitempty"` // 可选附加信息
}

func NewError(code ErrorCode, message string, data ...string) *Error {
	err := &Error{
		Code:    code,
		Message: message,
	}
	if len(data) > 0 {
		err.Data = data[0]
	}
	return err
}
func (r *Response) Success(id string, data interface{}) {
	r.JSONRPC = "2.0"
	r.ID = id
	r.Result = data
	r.Error = nil
}

func (r *Response) Errorf(id string, err *Error) {
	r.JSONRPC = "2.0"
	r.ID = id
	r.Error = err
	r.Result = nil
}
