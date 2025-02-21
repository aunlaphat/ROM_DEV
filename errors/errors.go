// Http Error Wrapper Package เก็บข้อผิดพลาดที่ใช้ในระบบ
package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NotFoundError(format string, a ...interface{}) error {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf(format, a...), 
	}
}

func UnexpectedError() error {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: "unexpected error",
	}
}

func ValidationError(format string, a ...interface{}) error {
	return &AppError{
		Code:    http.StatusUnprocessableEntity,
		Message: fmt.Sprintf(format, a...),
	}
}

func UnauthorizedError(format string, a ...interface{}) error {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: fmt.Sprintf(format, a...),
	}
}

func BadRequestError(format string, a ...interface{}) error {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf(format, a...),
	}
}

func InternalError(format string, a ...interface{}) error {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: fmt.Sprintf(format, a...),
	}
}

func ConflictError(format string, a ...interface{}) error {
	return &AppError{
		Code:    http.StatusConflict,
		Message: fmt.Sprintf(format, a...),
	}
}

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// อันเก่า
// func (e AppError) Error() string {
// 	return e.Message
// }

// func NotFoundError(message string) error {
// 	return AppError{
// 		Code:    http.StatusNotFound,
// 		Message: fmt.Sprintf("%v not found", message),
// 	}
// }

// func UnexpectedError() error {
// 	return AppError{
// 		Code:    http.StatusInternalServerError,
// 		Message: "unexpected error",
// 	}
// }

// func ValidationError(message string) error {
// 	return AppError{
// 		Code:    http.StatusUnprocessableEntity,
// 		Message: message,
// 	}
// }

// func UnauthorizedError(message string) error {
// 	return AppError{
// 		Code:    http.StatusUnauthorized,
// 		Message: message,
// 	}
// }

// func BadRequestError(message string) error {
// 	return AppError{
// 		Code:    http.StatusBadRequest,
// 		Message: message,
// 	}
// }

// func InternalError(message string) error {
// 	return AppError{
// 		Code:    http.StatusInternalServerError,
// 		Message: message,
// 	}
// }

// func ErrorHandler(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		defer func() {
// 			if rec := recover(); rec != nil {
// 				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 			}
// 		}()
// 		next.ServeHTTP(w, r)
// 	})
// }