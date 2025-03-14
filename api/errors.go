package api

import (
	Errors "boilerplate-back-go-2411/errors"
	"errors"
	"fmt"
	"net/http"

	// "boilerplate-back-go-2411/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func handleValidationError(c *gin.Context, err error) {
	var errorMessages []string

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			errorMessage := fmt.Sprintf("❌ Field '%s' is invalid: %s", fieldErr.Field(), fieldErr.Tag())
			errorMessages = append(errorMessages, errorMessage)
		}
	} else {
		errorMessages = append(errorMessages, err.Error())
	}

	handleResponse(c, false, "⚠️ Invalid request body", errorMessages, http.StatusBadRequest)
}

func handleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// 🔹 ตรวจสอบว่า Error มาจาก Service Layer
	var appErr *Errors.AppError
	if errors.As(err, &appErr) {
		handleResponse(c, false, appErr.Message, nil, appErr.Code)
		return
	}

	// // 🔹 ตรวจสอบว่า Error มาจาก Service Layer หรือไม่
	// if appErr, ok := err.(*Errors.AppError); ok {
	// 	handleResponse(c, false, "⚠️ Service error", appErr.Message, appErr.Code)
	// 	return
	// }

	handleResponse(c, false, "🔥 Internal server error", err.Error(), http.StatusInternalServerError)
}
