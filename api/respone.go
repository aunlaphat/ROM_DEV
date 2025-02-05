package api

import "github.com/gin-gonic/gin"

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func handleResponse(c *gin.Context, success bool, message string, data interface{}, status int) {
	c.JSON(status, Response{
		Success: success,
		Message: message,
		Data:    data,
	})
}
