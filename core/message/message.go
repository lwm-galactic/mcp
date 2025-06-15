package message

type Request struct {
	Method  string                 `json:"method"` // 如 getResource、invokeTool
	Params  map[string]interface{} `json:"params"`
	Session string                 `json:"session,omitempty"`
}

type Response struct {
	Status string      `json:"status"` // success / error
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}
