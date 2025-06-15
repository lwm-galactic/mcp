package handler

import "github.com/lwm-galactic/mcp/context"

type ResourceHandler func(ctx *context.ResourceContext) ([]byte, error)
