package api

import (
	Errors "boilerplate-backend-go/errors"
	"errors"
	"fmt"
	"net/http"

	"boilerplate-backend-go/errors"

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

	if appErr, ok := err.(errors.AppError); ok {
		handleResponse(c, false, "⚠️ Service error", appErr.Message, appErr.Code)
		return
	}

	handleResponse(c, false, "🔥 Internal server error", err.Error(), http.StatusInternalServerError)
}
