package errors

import (
	"fmt"
	"net/http"
)

type ApiError struct {
	Code    int
	Message string
}

func (e *ApiError) Error() string {
	return e.Message
}

func BadRequestError(format string, args ...any) *ApiError {
	return newApiError(http.StatusBadRequest, fmt.Sprintf(format, args...))
}

func InternalServerError(format string, args ...any) *ApiError {
	return newApiError(http.StatusInternalServerError, fmt.Sprintf(format, args...))
}

func newApiError(code int, message string) *ApiError {
	return &ApiError{
		Code:    code,
		Message: message,
	}
}
