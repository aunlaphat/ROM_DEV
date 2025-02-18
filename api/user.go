package api

import (
	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ✅ **ตั้งค่า API Route**
func (app *Application) UserRoute(apiRouter *gin.RouterGroup) {
	users := apiRouter.Group("/manage-users")

	// 🔹 Protected API (ต้องใช้ JWT)
	users.Use(middleware.JWTMiddleware(app.TokenAuth))

	users.GET("/", app.GetUsers)
	users.GET("/:userID", app.GetUser)
	users.POST("/add", app.AddUser)
	users.PATCH("/edit/:userID", app.EditUser)
	users.DELETE("/delete/:userID", app.DeleteUser)
}

// ✅ **1️⃣ GetUser - ดึงข้อมูลผู้ใช้**
// @Summary Get user details
// @Description Retrieve details of a specific user
// @Tags User Management
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {object} Response{data=response.UserResponse}
// @Failure 404 {object} Response
// @Router /manage-users/{userID} [get]
func (app *Application) GetUser(c *gin.Context) {
	userID := c.Param("userID")

	user, err := app.Service.User.GetUser(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}

	handleResponse(c, true, "⭐ User retrieved successfully ⭐", user, http.StatusOK)
}

// ✅ **2️⃣ GetUsers - ดึงรายชื่อผู้ใช้ทั้งหมด**
// @Summary Get list of users
// @Description Retrieve user data filtered by isActive, with pagination
// @Tags User Management
// @Accept json
// @Produce json
// @Param isActive query bool false "Filter by Active Status (true/false)"
// @Param limit query int false "Limit (default 100)"
// @Param offset query int false "Offset (default 0)"
// @Success 200 {object} Response{data=[]response.UserResponse}
// @Failure 400 {object} Response
// @Router /manage-users [get]
func (app *Application) GetUsers(c *gin.Context) {
	isActiveQuery := c.Query("isActive")

	var isActive bool
	if isActiveQuery != "" {
		parsedBool, err := strconv.ParseBool(isActiveQuery)
		if err != nil {
			handleResponse(c, false, "⚠️ Invalid isActive parameter (must be true/false)", nil, http.StatusBadRequest)
			return
		}
		isActive = parsedBool
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if err != nil || limit <= 0 {
		handleResponse(c, false, "⚠️ Invalid limit parameter", nil, http.StatusBadRequest)
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		handleResponse(c, false, "⚠️ Invalid offset parameter", nil, http.StatusBadRequest)
		return
	}

	users, err := app.Service.User.GetUsers(c.Request.Context(), isActive, limit, offset)
	if err != nil {
		handleError(c, err)
		return
	}

	handleResponse(c, true, "⭐ Users retrieved successfully ⭐", users, http.StatusOK)
}

// ✅ **3️⃣ AddUser - เพิ่มผู้ใช้ใหม่**
// @Summary Add a new user
// @Description Add a user with role assignment
// @Tags User Management
// @Accept json
// @Produce json
// @Param request body request.AddUserRequest true "User details"
// @Success 201 {object} Response{data=response.AddUserResponse}
// @Failure 400 {object} Response
// @Router /manage-users/add [post]
func (app *Application) AddUser(c *gin.Context) {
	var req request.AddUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleValidationError(c, err)
		return
	}

	adminID := c.MustGet("UserID").(string)
	adminRoleID := c.MustGet("RoleID").(int)

	newUser, err := app.Service.User.AddUser(c.Request.Context(), req, adminID, adminRoleID)
	if err != nil {
		handleError(c, err)
		return
	}

	handleResponse(c, true, "⭐ User added successfully ⭐", newUser, http.StatusCreated)
}

// ✅ **4️⃣ EditUser - แก้ไขข้อมูลผู้ใช้**
// @Summary Edit user details
// @Description Update role and warehouse of a user
// @Tags User Management
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Param request body request.EditUserRequest true "Updated user details"
// @Success 200 {object} Response{data=response.EditUserResponse}
// @Failure 400 {object} Response
// @Router /manage-users/edit/{userID} [patch]
func (app *Application) EditUser(c *gin.Context) {
	userID := c.Param("userID") // ดึง `userID` จาก Path Parameter

	var req request.EditUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleValidationError(c, err)
		return
	}

	// 🟢 ตรวจสอบว่า `userID` ที่รับจาก API ต้องตรงกับ `req.UserID` (กันข้อผิดพลาดจาก frontend)
	if req.UserID != "" && req.UserID != userID {
		handleResponse(c, false, "⚠️ User ID in request body does not match path parameter", nil, http.StatusBadRequest)
		return
	}

	// 🔹 ดึงข้อมูลผู้ใช้ที่ส่งคำขอแก้ไข
	adminID := c.MustGet("UserID").(string)
	adminRoleID := c.MustGet("RoleID").(int)

	// 🟢 **เรียกใช้ Service Layer**
	updatedUser, err := app.Service.User.EditUser(c.Request.Context(), req, adminID, adminRoleID)
	if err != nil {
		handleError(c, err)
		return
	}

	// ✅ **Response หากแก้ไขสำเร็จ**
	handleResponse(c, true, "⭐ User edited successfully ⭐", updatedUser, http.StatusOK)
}

// ✅ **5️⃣ DeleteUser - ลบผู้ใช้ (Soft Delete)**
// @Summary Delete a user (Soft Delete)
// @Description Remove user from the system but keep data in the database
// @Tags User Management
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /manage-users/delete/{userID} [delete]
func (app *Application) DeleteUser(c *gin.Context) {
	userID := c.Param("userID")
	adminID := c.MustGet("UserID").(string)
	adminRoleID := c.MustGet("RoleID").(int)

	err := app.Service.User.DeleteUser(c.Request.Context(), userID, adminID, adminRoleID)
	if err != nil {
		handleError(c, err)
		return
	}

	handleResponse(c, true, "⭐ User deleted successfully ⭐", nil, http.StatusOK)
}
