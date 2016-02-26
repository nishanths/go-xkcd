package xkcd

import (
	"fmt"
	"net/http"
)

// StatusError specifies the status code and status text
// for error responses from the xkcd API endpoint.
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

// Error returns a string representation of the StatusError.
func (e StatusError) Error() string {
	return fmt.Sprintf("%d: %s", e.StatusCode, e.StatusText)
}
