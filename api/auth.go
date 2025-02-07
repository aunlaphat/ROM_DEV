package api

import (
	"context"
	"net/http"

	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/dto/response"
	"boilerplate-backend-go/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 📌 กำหนดเส้นทางสำหรับ Authentication API
func (app *Application) AuthRoute(apiRouter *gin.RouterGroup) {
	auth := apiRouter.Group("/auth")
	auth.POST("/login", app.Login)              // Standard Login
	auth.POST("/login-lark", app.LoginFromLark) // Login via Lark

	// Routes requiring JWT authentication
	auth.Use(middleware.JWTAuthMiddleware(app.TokenAuth))
	auth.GET("/", app.CheckAuthen)   // Check authentication status
	auth.POST("/logout", app.Logout) // Logout
}

// Generate JWT token from user claims
func (app *Application) GenerateToken(user response.User) string {
	claims := map[string]interface{}{
		"userID":     user.UserID,
		"userName":   user.UserName,
		"roleID":     user.RoleID,
		"fullNameTH": user.FullNameTH,
		"nickName":   user.NickName,
		"department": user.DepartmentNo,
		"platform":   user.Platform,
	}

	_, tokenString, _ := app.TokenAuth.Encode(claims)
	return tokenString
}

// @Summary User Login
// @Description Authenticates user credentials and generates a JWT token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param login-request body request.LoginWeb true "User login credentials"
// @Success 200 {object} response.User "JWT token"
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
	user, err := app.Service.User.Login(ctx, req)
	if err != nil {
		app.Logger.Warn("⚠️ Login failed", zap.String("username", req.UserName), zap.Error(err))
		handleResponse(c, false, "Invalid username or password", nil, http.StatusUnauthorized)
		return
	}

	token := app.GenerateToken(user)
	app.Logger.Info("✅ Login successful", zap.String("username", user.UserName))

	c.SetCookie("jwt", token, 4*3600, "/", "", false, true)

	handleResponse(c, true, "Login Success", token, http.StatusOK)
}

// @Summary User Lark Login
// @Description Authenticates user credentials from Lark and generates a JWT token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param login-request body request.LoginLark true "User login from Lark"
// @Success 200 {object} response.User "JWT token"
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
	user, err := app.Service.User.LoginLark(ctx, req)
	if err != nil {
		app.Logger.Warn("⚠️ Login from Lark failed", zap.String("username", req.UserName), zap.String("userID", req.UserID), zap.Error(err))
		handleResponse(c, false, "User not found", nil, http.StatusUnauthorized)
		return
	}

	token := app.GenerateToken(user)
	app.Logger.Info("✅ Lark login successful", zap.String("username", user.UserName))

	c.SetCookie("jwt", token, 4*3600, "/", "", false, true)

	handleResponse(c, true, "Lark Login Success", token, http.StatusOK)
}

// @Summary User Logout
// @Description Logs the user out by removing the JWT token from the cookie.
// @Tags Auth
// @Success 200 {object} api.Response "Logout successful"
// @Router /auth/logout [post]
func (app *Application) Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "/", "", false, true)
	app.Logger.Info("✅ User logged out successfully")
	handleResponse(c, true, "Logout successful", nil, http.StatusOK)
}

// @Summary Check Authentication
// @Description Validates if the JWT token is valid and retrieves user claims.
// @Tags Auth
// @Success 200 {object} api.Response "Authenticated user details"
// @Failure 401 {object} api.Response "Unauthorized"
// @Router /auth [get]
func (app *Application) CheckAuthen(c *gin.Context) {
	source, _ := c.Get("jwt_source")

	var claims interface{}
	if source == "header" {
		claims, _ = c.Get("jwt_claims_header")
	} else if source == "cookie" {
		claims, _ = c.Get("jwt_claims_cookie")
	}

	//fmt.Printf("JWT Source: %s, Claims: %+v\n", source, claims)

	if claims == nil {
		handleResponse(c, false, "Unauthorized - No claims found", nil, http.StatusUnauthorized)
		return
	}

	handleResponse(c, true, "User authenticated", gin.H{
		"source": source,
		"claims": claims,
	}, http.StatusOK)
}
