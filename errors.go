package xkcd

import (
	"fmt"
	"net/http"
)

type StatusError struct {
	StatusCode int
	StatusText string
}

func newStatusError(code int) StatusError {
	return StatusError{
		StatusCode: code,
		StatusText: http.StatusText(code),
	}
}

func (e StatusError) Error() string {
	return fmt.Sprintf("%d: %s", e.StatusCode, e.StatusText)
}
