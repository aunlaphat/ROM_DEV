package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (app *Application) UserRoute(apiRouter *gin.RouterGroup) {
	user := apiRouter.Group("/user")
	user.GET("/:username", app.GetUser)
}

// @Summary Get User Credentials
// @Description Get user credentials by userName
// @Tags User
// @Produce json
// @Param username path string true "UserName"
// @Success 200 {object} response.UserRole "User credentials"
// @Failure 404 {object} gin.H "User not found"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /user/{username} [get]
func (app *Application) GetUser(c *gin.Context) {
	username := c.Param("username")

	user, err := app.Service.User.GetUser(context.Background(), username)
	if err != nil {
		if err.Error() == "user not found" {
			app.Logger.Warn("⚠️ User not found", zap.String("username", username))
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		app.Logger.Error("❌ Failed to get user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, user)
}
