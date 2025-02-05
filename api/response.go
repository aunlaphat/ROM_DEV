package api

import "github.com/gin-gonic/gin"

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// handleResponse is a helper function to send JSON responses
func handleResponse(c *gin.Context, success bool, message string, data interface{}, statusCode int) {
	c.JSON(statusCode, gin.H{
		"success": success,
		"message": message,
		"data":    data,
	})
}
