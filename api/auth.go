package api

import (
	req "boilerplate-backend-go/dto/request"
	res "boilerplate-backend-go/dto/response"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
)

// กำหนดเส้นทางสำหรับการตรวจสอบสิทธิ์
func (app *Application) AuthRoute(apiRouter *chi.Mux) {
	apiRouter.Route("/auth", func(r chi.Router) {
		r.Post("/login", app.Login)              // สำหรับ login
		r.Post("/login-lark", app.LoginFromLark) // สำหรับ login ผ่าน Lark
		r.Group(func(router chi.Router) {
			router.Use(jwtauth.Verifier(app.TokenAuth)) // middleware ตรวจสอบ token
			router.Use(jwtauth.Authenticator)           // middleware ยืนยันตัวตน
			router.Get("/", app.CheckAuthen)            // ตรวจสอบสถานะการ authentication
			router.Post("/logout", app.Logout)          // สำหรับ logout
		})
	})
}

var contentType = "content-type"
var appJson = "application/json"

// สร้าง JWT token โดยใช้ข้อมูลผู้ใช้
func (app *Application) GenerateToken(tokenData res.Login) string {
	// สร้าง claims (ข้อมูลที่จะเก็บใน token)
	data := map[string]interface{}{
		"userID":     tokenData.UserID,
		"userName":   tokenData.UserName,
		"roleID":     tokenData.RoleID,
		"fullNameTH": tokenData.FullNameTH,
		"nickName":   tokenData.NickName,
		"department": tokenData.DepartmentNo,
		"platform":   tokenData.Platform,
	}
	// สร้างและเข้ารหัส token
	_, tokenString, _ := app.TokenAuth.Encode(data)
	return tokenString
}

// @Summary User Login
// @Description Handles user login requests and generates a token for the authenticated user.
// @ID user-login
// @Tags Auth
// @Accept json
// @Produce json
// @Param login-request body request.LoginWeb true "User login credentials in JSON format"
// @Success 200 {object} Response{result=string} "JWT token"
// @Failure 400 {object} Response "Bad Request"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /auth/login [post]
// จัดการการเข้าสู่ระบบของผู้ใช้และสร้าง JWT token
func (app *Application) Login(w http.ResponseWriter, r *http.Request) {
	// 1. ตรวจสอบ content-type
	if r.Header.Get(contentType) != appJson {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// 2. รับข้อมูล login จาก request body
	req := req.LoginWeb{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(w, err)
		return
	}

	// 3. เรียกใช้ service เพื่อตรวจสอบ credentials
	user, err := app.Service.User.Login(req)
	if err != nil {
		handleError(w, err)
		return
	}
	tokenData := res.Login{
		UserID:       user.UserID,
		UserName:     user.UserName,
		RoleID:       user.RoleID,
		FullNameTH:   user.FullNameTH,
		NickName:     user.NickName,
		DepartmentNo: user.DepartmentNo,
		Platform:     user.Platform,
	}
	//fmt.Println("token data", tokenData)

	// 4. สร้าง JWT token จากข้อมูลผู้ใช้ (claims) -> func GenerateToken
	token := app.GenerateToken(tokenData)
	fmt.Println("token: ", token)

	// 5. ตั้งค่า cookie ที่มี token
	http.SetCookie(w, &http.Cookie{
		HttpOnly: false,
		Expires:  time.Now().Add(4 * time.Hour), //4 hours life
		SameSite: http.SameSiteLaxMode,
		Name:     "jwt", // Must be named "jwt" or else the token cannot be searched for by jwtauth.Verifier.
		Value:    token,
		Path:     "/",
	})
	// 6. ส่ง response กลับ
	handleResponse(w, true, "login Success", token, http.StatusOK)
}

// @Summary User Lark Login
// @Description Handles user login requests and generates a token for the authenticated user.
// @ID user-login-lark
// @Tags Auth
// @Accept json
// @Produce json
// @Param Login-request-lark body request.LoginLark true "User login from lark credentials from Lark in JSON format"
// @Success 200 {object} Response{result=string} "JWT token"
// @Failure 400 {object} Response "Bad Request"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /auth/login-lark [post]
// จัดการการเข้าสู่ระบบจาก Lark และสร้าง JWT token
func (app *Application) LoginFromLark(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("content-type") != "application/json" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	req := req.LoginLark{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(w, err)
		return
	}

	user, err := app.Service.User.LoginLark(req)
	if err != nil {
		handleError(w, err)
		return
	}
	tokenData := res.Login{
		UserID:       user.UserID,
		UserName:     user.UserName,
		RoleID:       user.RoleID,
		FullNameTH:   user.FullNameTH,
		NickName:     user.NickName,
		DepartmentNo: user.DepartmentNo,
		Platform:     user.Platform,
	}
	//fmt.Println("token data", tokenData)

	token := app.GenerateToken(tokenData)
	fmt.Println("token: ", token)

	http.SetCookie(w, &http.Cookie{
		HttpOnly: false,
		Expires:  time.Now().Add(4 * time.Hour), //4 hours life
		SameSite: http.SameSiteLaxMode,
		Name:     "jwt", // Must be named "jwt" or else the token cannot be searched for by jwtauth.Verifier.
		Value:    token,
		Path:     "/",
	})

	handleResponse(w, true, "login's Lark Success", token, http.StatusOK)
	// w.WriteHeader(http.StatusOK)

	// w.Header().Set("content-type", "application/json")
	// json.NewEncoder(w).Encode(user)
}

// @Summary User Logout
// @Description Logs out the user by deleting the JWT token.
// @ID user-logout
// @Tags Auth
// @Success 200 {object} Response{result=string} "Logout successful"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /auth/logout [post]
// จัดการการออกจากระบบของผู้ใช้โดยการลบ JWT token
func (app *Application) Logout(w http.ResponseWriter, r *http.Request) {
	// ลบ cookie โดยการตั้ง MaxAge เป็นค่าลบ
	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		MaxAge:   -1, // ลบ cookie
		SameSite: http.SameSiteLaxMode,
		Name:     "jwt",
		Value:    "",
	})
}

// @Summary Check Authentication
// @Description A test endpoint to check if the user is authenticated and to demonstrate Swagger documentation.
// @ID check-authentication
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} Response{result=map[string]interface{}} "Authenticated user details"
// @Failure 401 {object} Response "Unauthorized"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /auth [get]
// ตรวจสอบว่าผู้ใช้ได้รับการตรวจสอบสิทธิ์แล้วหรือไม่
func (app *Application) CheckAuthen(w http.ResponseWriter, r *http.Request) {
	// ดึง claims จาก context (ถูกเพิ่มโดย middleware)
	_, claims, _ := jwtauth.FromContext(r.Context())

	// ส่ง claims กลับเพื่อแสดงข้อมูลผู้ใช้
	handleResponse(w, true, "Checked", claims, http.StatusOK)
}
