package models

import "fmt"

type ErrorHandler struct {
	Message    interface{}
	StatusCode int
}

func (e *ErrorHandler) Error() string {
	return fmt.Sprintf("Error: %v StatusCode: %v", e.Message, e.StatusCode)
}
