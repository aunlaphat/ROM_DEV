package api

import (
	"context"
	"fmt"
	"net/http"

	"boilerplate-back-go-2411/dto/request"
	"boilerplate-back-go-2411/dto/response"
	"boilerplate-back-go-2411/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (app *Application) AuthRoute(apiRouter *gin.RouterGroup) {
	auth := apiRouter.Group("/auth")
	auth.POST("/login", app.Login)
	auth.POST("/login-lark", app.LoginFromLark)
	auth.POST("/logout", app.Logout)

	protected := auth.Group("/")
	protected.Use(middleware.JWTMiddleware(app.TokenAuth))
	protected.GET("/", app.CheckAuthen)
}

func (app *Application) GenerateToken(tokenData response.Login) string {
	claims := map[string]interface{}{
		"userID":       tokenData.UserID,
		"userName":     tokenData.UserName,
		"fullNameTH":   tokenData.FullNameTH,
		"nickName":     tokenData.NickName,
		"roleID":       tokenData.RoleID,
		"roleName":     tokenData.RoleName,
		"departmentNo": tokenData.DepartmentNo,
		"platform":     tokenData.Platform,
	}

	_, tokenString, _ := app.TokenAuth.Encode(claims)

	fmt.Println("🔑 Debug JWT:", tokenString)

	return tokenString
}

// @Summary User Login
// @Description Authenticates user credentials and generates a JWT token.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param login-request body request.LoginWeb true "User login credentials"
// @Success 200 {object} response.Login "JWT token"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /auth/login [post]
func (app *Application) Login(c *gin.Context) {
	var req request.LoginWeb
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, false, "Invalid request payload", nil, http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	user, err := app.Service.Auth.Login(ctx, req)
	if err != nil {
		app.Logger.Warn("⚠️ Login failed", zap.String("username", req.UserName), zap.Error(err))
		handleResponse(c, false, "Invalid username or password", nil, http.StatusUnauthorized)
		return
	}

	token := app.GenerateToken(user)
	app.Logger.Info("✅ Login successful", zap.String("username", user.UserName))

	if token == "" {
		handleResponse(c, false, "Failed to generate token", nil, http.StatusInternalServerError)
		return
	}

	c.SetCookie("jwt", token, 4*3600, "/", "", false, true) // Secure: false for localhost

	handleResponse(c, true, "Login Success", token, http.StatusOK)
}

// @Summary User Lark Login
// @Description Authenticates user credentials from Lark and generates a JWT token.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param login-request body request.LoginLark true "User login from Lark"
// @Success 200 {object} response.Login "JWT token"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /auth/login-lark [post]
func (app *Application) LoginFromLark(c *gin.Context) {
	var req request.LoginLark
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, false, "Invalid request payload", nil, http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	user, err := app.Service.Auth.LoginLark(ctx, req)
	if err != nil {
		app.Logger.Warn("⚠️ Login from Lark failed", zap.String("username", req.UserName), zap.String("userID", req.UserID), zap.Error(err))
		handleResponse(c, false, "User not found", nil, http.StatusUnauthorized)
		return
	}

	token := app.GenerateToken(user)
	app.Logger.Info("✅ Lark login successful", zap.String("username", user.UserName))

	c.SetCookie("jwt", token, 4*3600, "/", "localhost", false, true)

	handleResponse(c, true, "Lark Login Success", token, http.StatusOK)
}

// @Summary User Logout
// @Description Logs the user out by removing the JWT token from the cookie.
// @Tags Authentication
// @Success 200 {object} api.Response "Logout successful"
// @Router /auth/logout [post]
func (app *Application) Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "/", "", false, true)
	app.Logger.Info("✅ User logged out successfully")
	handleResponse(c, true, "Logout successful", nil, http.StatusOK)
}

// @Summary Check Authentication
// @Description Validates if the JWT token is valid and retrieves user claims.
// @Tags Authentication
// @Success 200 {object} api.Response "Authenticated user details"
// @Failure 401 {object} api.Response "Unauthorized"
// @Router /auth [get]
func (app *Application) CheckAuthen(c *gin.Context) {
	// ✅ ดึง JWT Source (Header หรือ Cookie)
	source, _ := c.Get("jwt_source")
	claims := middleware.GetJWTClaims(c)

	if claims == nil {
		handleResponse(c, false, "Unauthorized - No claims found", nil, http.StatusUnauthorized)
		return
	}

	app.Logger.Info("✅ User authenticated", zap.Any("claims", claims))

	handleResponse(c, true, "User authenticated", gin.H{
		"source": source,
		"user": gin.H{
			"userID":       claims["userID"],
			"userName":     claims["userName"],
			"fullNameTH":   claims["fullNameTH"],
			"nickName":     claims["nickName"],
			"roleID":       claims["roleID"],
			"roleName":     claims["roleName"],
			"departmentNo": claims["departmentNo"],
			"platform":     claims["platform"],
		},
	}, http.StatusOK)
}
