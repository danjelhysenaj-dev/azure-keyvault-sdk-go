package errors

import (
	"fmt"
	"net/http"
)

const (
	ErrCodeInternalServerError = "INTERNAL_SERVER_ERROR"
	ErrCodeNotFound            = "NOT_FOUND"
	ErrCodeUnathorized         = "UNAUTHORIZED"
	ErrCodeInsufficientAccess  = "INSUFFICIENT_ACCESS"
)

// Errors represents a list of Error.
type Errors struct {
	Errors []Error `json:"errors"`
}

// Error represents the error information.
type Error struct {
	Code    string `json:"code"`
	Status  int    `json:"status"`
	Message string `json:"message"`
	TraceId string `json:"traceId"`
}

// Error return the string formatted Error.
func (e *Error) Error() string {
	msg := fmt.Sprintf("Code: %s, Status: %d, Message: %s, TraceId: %s", e.Code, e.Status, e.Message, e.TraceId)
	return msg
}

func InternalServerError(message string) *Error {
	return &Error{
		Code:    ErrCodeInternalServerError,
		Status:  http.StatusInternalServerError,
		Message: message,
		TraceId: "",
	}
}

func InternalServerErrorf(format string, a ...any) *Error {
	return &Error{
		Code:    ErrCodeInternalServerError,
		Status:  http.StatusInternalServerError,
		Message: fmt.Sprintf(format, a...),
		TraceId: "",
	}
}

func NotFoundError(message string) *Error {
	return &Error{
		Code:    ErrCodeNotFound,
		Status:  http.StatusNotFound,
		Message: message,
		TraceId: "",
	}
}

func UnAuthorizedError(message string) *Error {
	return &Error{
		Code:    ErrCodeUnathorized,
		Status:  http.StatusUnauthorized,
		Message: message,
		TraceId: "",
	}
}

func ForbiddenError(message string) *Error {
	return &Error{
		Code:    ErrCodeInsufficientAccess,
		Status:  http.StatusForbidden,
		Message: message,
		TraceId: "",
	}
}
