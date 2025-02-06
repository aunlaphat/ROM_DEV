package api

import (
	"github.com/gin-gonic/gin"
)

// 📌 กำหนดเส้นทาง API สำหรับ User
func (app *Application) UserRoute(apiRouter *gin.RouterGroup) {
	/* user := apiRouter.Group("/user")

	user.POST("/get-user", app.GetUser)
	user.POST("/get-user-with-permission", app.GetUserWithPermission) */
}

/*
// GetUser godoc
// @Summary Get user by userid and username
// @Description Retrieve the details of a user by their userid and username
// @ID get-user
// @Tags User
// @Accept json
// @Produce json
// @Param Login body request.LoginLark true "User login credentials in JSON format"
// @Success 200 {object} api.Response{data=response.Login} "User retrieved successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 404 {object} api.Response "User not found"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /user/get-user [post]
func (app *Application) GetUser(c *gin.Context) {
	// ✅ รับข้อมูลจาก Request Body
	var req request.LoginLark
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, false, "Invalid request payload", nil, http.StatusBadRequest)
		return
	}

	// ✅ ตรวจสอบค่า UserID และ UserName
	if req.UserID == "" || req.UserName == "" {
		handleResponse(c, false, "UserID and UserName are required", nil, http.StatusBadRequest)
		return
	}

	// ✅ Log ก่อนเรียกใช้งาน Service
	app.Logger.Info("🔍 Searching for user",
		zap.String("UserID", req.UserID),
		zap.String("UserName", req.UserName),
	)

	// ✅ ค้นหาข้อมูล User
	user, err := app.Service.User.GetUser(c, req)
	if err != nil {
		app.Logger.Warn("⚠️ User not found",
			zap.String("UserID", req.UserID),
			zap.String("UserName", req.UserName),
		)
		handleResponse(c, false, "User not found", nil, http.StatusNotFound)
		return
	}

	// ✅ Logging & Response
	app.Logger.Info("✅ User retrieved successfully", zap.String("UserID", user.UserID))
	handleResponse(c, true, "🤹🏻 User retrieved successfully 🤹🏻", user, http.StatusOK)
}

// GetUserWithPermission godoc
// @Summary Get user with permissions by username and password
// @Description Retrieve the details of a user with permissions by their username and password
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
	// ✅ รับข้อมูลจาก Request Body
	var req request.LoginLark
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, false, "Invalid request payload", nil, http.StatusBadRequest)
		return
	}

	// ✅ ตรวจสอบค่า UserID และ UserName
	if req.UserID == "" || req.UserName == "" {
		handleResponse(c, false, "UserID and UserName are required", nil, http.StatusBadRequest)
		return
	}

	// ✅ Log ก่อนเรียกใช้งาน Service
	app.Logger.Info("🔍 Searching for user with permissions",
		zap.String("UserID", req.UserID),
		zap.String("UserName", req.UserName),
	)

	// ✅ ค้นหาข้อมูล User พร้อมสิทธิ์การใช้งาน
	user, err := app.Service.User.GetUserWithPermission(c, req)
	if err != nil {
		app.Logger.Warn("⚠️ User not found",
			zap.String("UserID", req.UserID),
			zap.String("UserName", req.UserName),
		)
		handleResponse(c, false, "User not found", nil, http.StatusNotFound)
		return
	}

	// ✅ Logging & Response
	app.Logger.Info("✅ User with permissions retrieved successfully", zap.String("UserID", user.UserID))
	handleResponse(c, true, "🤹🏻 User with permissions retrieved successfully 🤹🏻", user, http.StatusOK)
}
*/
