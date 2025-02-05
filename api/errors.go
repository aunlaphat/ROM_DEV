package api

import (
	"boilerplate-backend-go/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleError(c *gin.Context, err error) {
	switch e := err.(type) {
	case errors.AppError:
		c.JSON(e.Code, gin.H{"error": e.Message})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
