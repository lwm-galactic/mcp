package message

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

type Error struct {
	Code    ErrorCode   `json:"code"`           // 错误码
	Message string      `json:"message"`        // 错误描述
	Data    interface{} `json:"data,omitempty"` // 可选附加信息
}

func newJSONRPCError(code ErrorCode, message string, data ...string) *Error {
	err := &Error{
		Code:    code,
		Message: message,
	}
	if len(data) > 0 {
		err.Data = data[0]
	}
	return err
}
