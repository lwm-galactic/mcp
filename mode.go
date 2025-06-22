package zeno

import (
	"io"
	"os"
)

type Mode string

var (
	mcpMode = DebugMode
)

const (
	DebugMode   Mode = "debug"
	ReleaseMode Mode = "release"
	TestMode    Mode = "test"
)

var DefaultWriter io.Writer = os.Stdout

var DefaultErrorWriter io.Writer = os.Stderr

func SetMode(value Mode) {
	mcpMode = value
}

func GetMode() Mode {
	return mcpMode
}
