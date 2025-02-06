package api

import (
	"context"
	"net/http"

	"boilerplate-backend-go/dto/request"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (app *Application) UserRoute(apiRouter *gin.RouterGroup) {
	user := apiRouter.Group("/user")

	user.POST("/get-user", app.GetUser)
	user.POST("/get-user-with-permission", app.GetUserWithPermission)
}

// GetUser godoc
// @Summary Get user by username
// @Description Retrieve user details by username
// @ID get-user
// @Tags User
// @Accept json
// @Produce json
// @Param Login body request.LoginWeb true "User login credentials in JSON format"
// @Success 200 {object} api.Response{data=response.User} "User retrieved successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 404 {object} api.Response "User not found"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /user/get-user [post]
func (app *Application) GetUser(c *gin.Context) {
	var req request.LoginWeb
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, false, "Invalid request payload", nil, http.StatusBadRequest)
		return
	}

	if req.UserName == "" {
		handleResponse(c, false, "UserName is required", nil, http.StatusBadRequest)
		return
	}

	app.Logger.Info("🔍 Searching for user",
		zap.String("UserName", req.UserName),
	)

	user, err := app.Service.User.GetUser(context.Background(), req.UserName)
	if err != nil {
		app.Logger.Warn("⚠️ User not found",
			zap.String("UserName", req.UserName),
		)
		handleResponse(c, false, "User not found", nil, http.StatusNotFound)
		return
	}

	app.Logger.Info("✅ User retrieved successfully", zap.String("UserID", user.UserID))
	handleResponse(c, true, "🤹🏻 User retrieved successfully 🤹🏻", user, http.StatusOK)
}

// ✅ **GetUserWithPermission API**
// @Summary Get user with permissions
// @Description Retrieve user details along with role permissions
// @ID get-user-with-permission
// @Tags User
// @Accept json
// @Produce json
// @Param Login body request.LoginLark true "User login credentials in JSON format"
// @Success 200 {object} api.Response{data=response.UserPermission} "User with permissions retrieved successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 404 {object} api.Response "User not found"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /user/get-user-with-permission [post]
func (app *Application) GetUserWithPermission(c *gin.Context) {
	var req request.LoginLark
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, false, "Invalid request payload", nil, http.StatusBadRequest)
		return
	}

	if req.UserID == "" || req.UserName == "" {
		handleResponse(c, false, "UserID and UserName are required", nil, http.StatusBadRequest)
		return
	}

	app.Logger.Info("🔍 Searching for user with permissions",
		zap.String("UserID", req.UserID),
		zap.String("UserName", req.UserName),
	)

	userPermission, err := app.Service.User.GetUserWithPermission(context.Background(), req.UserID, req.UserName)
	if err != nil {
		app.Logger.Warn("⚠️ User not found",
			zap.String("UserID", req.UserID),
			zap.String("UserName", req.UserName),
		)
		handleResponse(c, false, "User not found", nil, http.StatusNotFound)
		return
	}

	app.Logger.Info("✅ User with permissions retrieved successfully", zap.String("UserID", userPermission.UserID))
	handleResponse(c, true, "🤹🏻 User with permissions retrieved successfully 🤹🏻", userPermission, http.StatusOK)
}
