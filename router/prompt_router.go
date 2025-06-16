package router

import (
	"fmt"
	"github.com/lwm-galactic/mcp/core/prompt"
)

type PromptRegistry struct {
	prompts map[string]prompt.Prompt
}

func NewPromptRegistry() *PromptRegistry {
	return &PromptRegistry{
		prompts: make(map[string]prompt.Prompt),
	}
}

// Register 添加一个提示到注册器中
func (r *PromptRegistry) Register(prompt prompt.Prompt) error {
	if _, exists := r.prompts[prompt.Name()]; exists {
		return fmt.Errorf("prompt already registered: %s", prompt.Name())
	}
	r.prompts[prompt.Name()] = prompt
	return nil
}

// Get 获取指定名称的提示
func (r *PromptRegistry) Get(name string) (prompt.Prompt, bool) {
	p, exists := r.prompts[name]
	return p, exists
}

// ListPrompts 返回所有已注册的提示信息
func (r *PromptRegistry) ListPrompts() []struct {
	Name        string
	Description string
	Arguments   []prompt.Argument
} {
	var result []struct {
		Name        string
		Description string
		Arguments   []prompt.Argument
	}
	for _, p := range r.prompts {
		result = append(result, struct {
			Name        string
			Description string
			Arguments   []prompt.Argument
		}{
			Name:        p.Name(),
			Description: p.Description(),
			Arguments:   p.Arguments(),
		})
	}
	return result
}

// Render 渲染指定名称的提示
func (r *PromptRegistry) Render(name string, args map[string]interface{}) ([]prompt.Message, error) {
	p, exists := r.prompts[name]
	if !exists {
		return nil, fmt.Errorf("prompt not found: %s", name)
	}
	return p.Render(args)
}
