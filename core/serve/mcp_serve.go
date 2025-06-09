package serve

import "github.com/lwm-galactic/mcp/core/tools"

type McpServe struct {
	Describe string
	manager  tools.ToolManager
}

func NewMcpServe(describe string) *McpServe {
	return &McpServe{
		Describe: describe,
	}
}

func (mcp *McpServe) SetToolManager(manager tools.ToolManager) {
	mcp.manager = manager
}

func (mcp *McpServe) SetResource() error {
	return nil
}

func (mcp *McpServe) Run() error {
	return nil
}
