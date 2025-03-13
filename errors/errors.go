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

// NotFoundError - ใช้สำหรับกรณีข้อมูลไม่พบ (404)
func NotFoundError(message string) error {
	return AppError{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf(format, a...), 
	}
}

// ConflictError - ใช้เมื่อข้อมูลซ้ำกัน (409 Conflict)
func ConflictError(message string) error {
	return AppError{
		Code:    http.StatusConflict,
		Message: fmt.Sprintf("%v : conflict", message),
	}
}

// ValidationError - ใช้เมื่อข้อมูลจากผู้ใช้ไม่ถูกต้อง (422 Unprocessable Entity)
func ValidationError(message string) error {
	return AppError{
		Code:    http.StatusUnprocessableEntity,
		Message: fmt.Sprintf(format, a...),
	}
}

// UnauthorizedError - ใช้เมื่อผู้ใช้ไม่มีสิทธิ์ใช้งาน (401)
func UnauthorizedError(message string) error {
	return AppError{
		Code:    http.StatusUnauthorized,
		Message: fmt.Sprintf(format, a...),
	}
}

// BadRequestError - ใช้สำหรับข้อมูลที่ไม่ถูกต้องจากฝั่ง Client (400)
func BadRequestError(message string) error {
	return AppError{
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf(format, a...),
	}
}

// InternalError - ใช้สำหรับข้อผิดพลาดที่ไม่คาดคิด (500)
func InternalError(message string) error {
	return AppError{
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
