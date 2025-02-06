package api

import (
	"context"
	"net/http"

	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/dto/response"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/jwtauth"
	"go.uber.org/zap"
)

// 📌 กำหนดเส้นทางสำหรับ Authentication API
func (app *Application) AuthRoute(apiRouter *gin.RouterGroup) {
	auth := apiRouter.Group("/auth")
	auth.POST("/login", app.Login)              // Login ปกติ
	auth.POST("/login-lark", app.LoginFromLark) // Login ผ่าน Lark

	// Routes ที่ต้องมี JWT
	auth.Use(jwtauth.Verifier(app.TokenAuth)) // middleware ตรวจสอบ token
	auth.Use(jwtauth.Authenticator)           // middleware ยืนยันตัวตน
	auth.GET("/", app.CheckAuthen)            // ตรวจสอบ Authentication
	auth.POST("/logout", app.Logout)          // Logout
}

// ✅ **GenerateToken()**
// ฟังก์ชันสร้าง JWT Token สำหรับผู้ใช้
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

	// สร้างและเข้ารหัส Token
	_, tokenString, _ := app.TokenAuth.Encode(claims)
	return tokenString
}

// ✅ **Login API**
// @Summary User Login
// @Description ตรวจสอบ credentials และออก JWT Token
// @Tags Auth
// @Accept json
// @Produce json
// @Param login-request body request.LoginWeb true "User login credentials in JSON format"
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

	// ✅ ตรวจสอบ username และ password
	ctx := context.Background()
	user, err := app.Service.User.Login(ctx, req)
	if err != nil {
		app.Logger.Warn("⚠️ Login failed", zap.String("username", req.UserName), zap.Error(err))
		handleResponse(c, false, "Invalid username or password", nil, http.StatusUnauthorized)
		return
	}

	// ✅ สร้าง JWT Token
	token := app.GenerateToken(user)
	app.Logger.Info("✅ Login successful", zap.String("username", user.UserName))

	// ✅ ตั้งค่า Cookie ให้ JWT
	c.SetCookie("jwt", token, 4*3600, "/", "", false, true) // 4 ชั่วโมง

	// ✅ ส่ง Response กลับ
	handleResponse(c, true, "Login Success", token, http.StatusOK)
}

// ✅ **Login ผ่าน Lark**
// @Summary User Lark Login
// @Description ตรวจสอบ Lark Credentials และออก JWT Token
// @Tags Auth
// @Accept json
// @Produce json
// @Param login-request body request.LoginLark true "User login from Lark in JSON format"
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

	// ✅ ตรวจสอบ username และ userID จาก Lark
	ctx := context.Background()
	user, err := app.Service.User.LoginLark(ctx, req)
	if err != nil {
		app.Logger.Warn("⚠️ Login from Lark failed", zap.String("username", req.UserName), zap.String("userID", req.UserID), zap.Error(err))
		handleResponse(c, false, "User not found", nil, http.StatusUnauthorized)
		return
	}

	// ✅ สร้าง JWT Token
	token := app.GenerateToken(user)
	app.Logger.Info("✅ Lark login successful", zap.String("username", user.UserName))

	// ✅ ตั้งค่า Cookie ให้ JWT
	c.SetCookie("jwt", token, 4*3600, "/", "", false, true)

	// ✅ ส่ง Response กลับ
	handleResponse(c, true, "Lark Login Success", token, http.StatusOK)
}

// ✅ **Logout API**
// @Summary User Logout
// @Description ลบ JWT Token ออกจาก Cookie
// @Tags Auth
// @Success 200 {object} api.Response "Logout successful"
// @Router /auth/logout [post]
func (app *Application) Logout(c *gin.Context) {
	// ✅ ลบ Cookie โดยการตั้ง MaxAge เป็น -1
	c.SetCookie("jwt", "", -1, "/", "", false, true)
	app.Logger.Info("✅ User logged out successfully")
	handleResponse(c, true, "Logout successful", nil, http.StatusOK)
}

// ✅ **Check Authentication API**
// @Summary Check Authentication
// @Description ตรวจสอบว่า JWT Token ถูกต้องหรือไม่
// @Tags Auth
// @Success 200 {object} api.Response "Authenticated user details"
// @Failure 401 {object} api.Response "Unauthorized"
// @Router /auth [get]
func (app *Application) CheckAuthen(c *gin.Context) {
	// ✅ ดึง claims จาก context (ถูกเพิ่มโดย middleware)
	_, claims, _ := jwtauth.FromContext(c.Request.Context())

	if claims == nil {
		handleResponse(c, false, "Unauthorized", nil, http.StatusUnauthorized)
		return
	}

	// ✅ ส่ง claims กลับ
	app.Logger.Info("✅ User authenticated", zap.Any("claims", claims))
	handleResponse(c, true, "User authenticated", claims, http.StatusOK)
}
