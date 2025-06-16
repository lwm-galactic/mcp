package prompt

type Prompt interface {
	Name() string
	Description() string
	Arguments() []Argument
	Render(args map[string]interface{}) ([]Message, error)
}
type Argument struct {
	Name        string
	Description string
	Required    bool
}

type Message struct {
	Role    string `json:"role"`    // "user", "assistant", "system"
	Content string `json:"content"` // 提示文本或 JSON 内容
}
