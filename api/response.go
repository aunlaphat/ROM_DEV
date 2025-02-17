package api

import "github.com/gin-gonic/gin"

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func handleResponse(c *gin.Context, success bool, message string, data interface{}, statusCode int) {
	c.JSON(statusCode, Response{
		Success: success,
		Message: message,
		Data:    data,
	})
}
