package api

import (
	request "boilerplate-backend-go/dto/request"
	response "boilerplate-backend-go/dto/response"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/jwtauth"
	"go.uber.org/zap"
)

// 📌 กำหนดเส้นทาง Authentication
func (app *Application) AuthRoute(apiRouter *gin.RouterGroup) {
	auth := apiRouter.Group("/auth")

	auth.POST("/login", app.Login)              // 🔹 Login
	auth.POST("/login-lark", app.LoginFromLark) // 🔹 Login ผ่าน Lark

	// 🔹 Protected Routes (ต้องใช้ JWT)
	auth.Use(func(c *gin.Context) {
		_, claims, err := jwtauth.FromContext(c.Request.Context())
		if err != nil {
			handleResponse(c, false, "❌ Unauthorized", nil, http.StatusUnauthorized)
			c.Abort()
			return
		}
		app.Logger.Info("🔑 JWT Claims", zap.Any("claims", claims))
		c.Set("jwt_claims", claims)
		c.Next()
	})

	auth.GET("/", app.CheckAuthen)   // 🔹 ตรวจสอบสิทธิ์
	auth.POST("/logout", app.Logout) // 🔹 Logout
}

// 📌 สร้าง JWT Token
func (app *Application) GenerateToken(tokenData response.Login) string {
	claims := map[string]interface{}{
		"userID":     tokenData.UserID,
		"userName":   tokenData.UserName,
		"roleID":     tokenData.RoleID,
		"fullNameTH": tokenData.FullNameTH,
		"nickName":   tokenData.NickName,
		"department": tokenData.DepartmentNo,
		"platform":   tokenData.Platform,
	}
	_, tokenString, _ := app.TokenAuth.Encode(claims)
	app.Logger.Info(fmt.Sprintf("🔑 JWT Claims: %+v", claims))
	app.Logger.Info("🔑 JWT Token: " + tokenString)

	return tokenString
}

// 📌 User Login
// @Summary User Login
// @Description Authenticate user and generate JWT token
// @ID user-login
// @Tags Auth
// @Accept json
// @Produce json
// @Param login-request body request.LoginWeb true "User credentials"
// @Success 200 {object} api.Response{data=string} "JWT token"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /auth/login [post]
func (app *Application) Login(c *gin.Context) {
	var req request.LoginWeb
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, false, "Invalid request", nil, http.StatusBadRequest)
		return
	}

	// 🔹 ตรวจสอบ User Credentials
	user, err := app.Service.User.Login(req)
	if err != nil {
		handleResponse(c, false, err.Error(), nil, http.StatusUnauthorized)
		return
	}

	// 🔹 สร้าง JWT Token
	token := app.GenerateToken(response.Login{
		UserID:       user.UserID,
		UserName:     user.UserName,
		RoleID:       user.RoleID,
		FullNameTH:   user.FullNameTH,
		NickName:     user.NickName,
		DepartmentNo: user.DepartmentNo,
		Platform:     user.Platform,
	})

	// 🔹 ตั้งค่า Cookie ที่มี Token
	c.SetCookie("jwt", token, 4*3600, "/", "", false, true) // 4 ชั่วโมง

	// 🔹 ส่ง Response กลับ
	handleResponse(c, true, "🟢 Login Success", token, http.StatusOK)
}

// 📌 User Login ผ่าน Lark
// @Summary User Lark Login
// @Description Authenticate user from Lark and generate JWT token
// @ID user-login-lark
// @Tags Auth
// @Accept json
// @Produce json
// @Param login-request-lark body request.LoginLark true "User credentials from Lark"
// @Success 200 {object} api.Response{data=string} "JWT token"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /auth/login-lark [post]
func (app *Application) LoginFromLark(c *gin.Context) {
	var req request.LoginLark
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, false, "Invalid request", nil, http.StatusBadRequest)
		return
	}

	user, err := app.Service.User.LoginLark(req)
	if err != nil {
		handleResponse(c, false, err.Error(), nil, http.StatusUnauthorized)
		return
	}

	token := app.GenerateToken(response.Login{
		UserID:       user.UserID,
		UserName:     user.UserName,
		RoleID:       user.RoleID,
		FullNameTH:   user.FullNameTH,
		NickName:     user.NickName,
		DepartmentNo: user.DepartmentNo,
		Platform:     user.Platform,
	})

	c.SetCookie("jwt", token, 4*3600, "/", "", false, true)

	handleResponse(c, true, "🟢 Login via Lark Success", token, http.StatusOK)
}

// 📌 User Logout
// @Summary User Logout
// @Description Logout user by deleting JWT token
// @ID user-logout
// @Tags Auth
// @Success 200 {object} api.Response "Logout successful"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /auth/logout [post]
func (app *Application) Logout(c *gin.Context) {
	// ลบ Cookie
	c.SetCookie("jwt", "", -1, "/", "", false, true)
	handleResponse(c, true, "🔴 Logout Success", nil, http.StatusOK)
}

// 📌 ตรวจสอบ Authentication
// @Summary Check Authentication
// @Description Check if the user is authenticated
// @ID check-authentication
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} api.Response{data=map[string]interface{}} "Authenticated user details"
// @Failure 401 {object} api.Response "Unauthorized"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /auth [get]
func (app *Application) CheckAuthen(c *gin.Context) {
	claims, exists := c.Get("jwt_claims")

	if !exists || claims == nil {
		handleResponse(c, false, "❌ Unauthorized", nil, http.StatusUnauthorized)
		return
	}

	claimsMap, ok := claims.(map[string]interface{})
	if !ok {
		app.Logger.Error("🚨 Invalid JWT Claims Format", zap.Any("claims", claims))
		handleResponse(c, false, "❌ Invalid Token Data", nil, http.StatusInternalServerError)
		return
	}

	handleResponse(c, true, "🟢 Authentication Checked 🟢", claimsMap, http.StatusOK)
}
