package api

import (
	"fmt"
	"net/http"

	"boilerplate-backend-go/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ✅ ฟังก์ชันจัดการ Validation Error (ใช้กับ JSON Binding)
func handleValidationError(c *gin.Context, err error) {
	var errorMessages []string

	// 🔹 ตรวจสอบว่าเป็น Validation Error หรือไม่
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			errorMessage := fmt.Sprintf("❌ Field '%s' is invalid: %s", fieldErr.Field(), fieldErr.Tag())
			errorMessages = append(errorMessages, errorMessage)
		}
	} else {
		// 🔹 หากเป็น Error อื่นที่ไม่ใช่ Validation Error
		errorMessages = append(errorMessages, err.Error())
	}

	// 🔹 ส่ง Response กลับไปพร้อมรายละเอียด Error
	handleResponse(c, false, "⚠️ Invalid request body", errorMessages, http.StatusBadRequest)
}

// ✅ ฟังก์ชันจัดการ Error ที่ส่งมาจาก Service Layer
func handleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// 🔹 ตรวจสอบว่า Error มาจาก Service Layer หรือไม่
	if appErr, ok := err.(errors.AppError); ok {
		handleResponse(c, false, "⚠️ Service error", appErr.Message, appErr.Code)
		return
	}

	// 🔹 หากเป็น Error อื่นที่ไม่สามารถระบุได้
	handleResponse(c, false, "🔥 Internal server error", err.Error(), http.StatusInternalServerError)
}
