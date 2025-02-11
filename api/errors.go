package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ✅ ฟังก์ชันจัดการ Error สำหรับ Validation (เช่น JSON Bind)
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

// ✅ ฟังก์ชันจัดการ Error ทั่วไป
func handleError(c *gin.Context, err error) {
	if err != nil {
		handleResponse(c, false, "🔥 Internal server error", err.Error(), http.StatusInternalServerError)
	}
}
