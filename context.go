package zeno

import (
	"net/http"
)

type Context struct {
	Request *http.Request
}
