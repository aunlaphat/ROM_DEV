// Http Error Wrapper Package
package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int
	Message string
}

func (e AppError) Error() string {
	return e.Message
}

func NotFoundError(message string) error {
	return AppError{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("%v not found", message),
	}
}

func UnexpectedError() error {
	return AppError{
		Code:    http.StatusInternalServerError,
		Message: "unexpected error",
	}
}

func ValidationError(message string) error {
	return AppError{
		Code:    http.StatusUnprocessableEntity,
		Message: message,
	}
}

func UnauthorizedError(message string) error {
	return AppError{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}
