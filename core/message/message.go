package message

import (
	"encoding/json"
	"fmt"
)

// Request 统一接收请求
type Request struct {
	ID      string                 `json:"id,omitempty"`      // 可选：用于追踪请求
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
	ID     string      `json:"id,omitempty"`    // 对应请求 ID
	Status string      `json:"status"`          // "success" / "error"
	Data   interface{} `json:"data,omitempty"`  // 返回数据
	Error  string      `json:"error,omitempty"` // 错误信息
}

func (r *Response) Success(data interface{}) {
	r.Status = "success"
	r.Data = data
	r.Error = ""
}

func (r *Response) Errorf(format string, args ...interface{}) {
	r.Status = "error"
	r.Data = nil
	r.Error = fmt.Sprintf(format, args...)
}
