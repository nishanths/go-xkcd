package xkcd

import (
	"fmt"
	"net/http"
)

type Error struct {
	StatusCode int
	StatusText string
}

func newError(code int) Error {
	return Error{
		StatusCode: code,
		StatusText: http.StatusText(code),
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("%d: %s", e.StatusCode, e.StatusText)
}
