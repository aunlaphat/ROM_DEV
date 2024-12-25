package api

import (
	req "boilerplate-backend-go/dto/request"
	res "boilerplate-backend-go/dto/response"
	"encoding/json"
	"fmt"
	"net/http"
//	"time"

	"github.com/go-chi/chi/v5"
	//"github.com/go-chi/jwtauth"
)

func (app *Application) UserTestRoute(apiRouter *chi.Mux) {

	apiRouter.Route("/test", func(r chi.Router) {
		r.Post("/login", app.LoginTest)
	 })
}

func (app *Application) GenerateTokenLogin(user res.UserInform) string {
	data := map[string]interface{}{
		"userID":   user.UserID,
		"userName": user.UserName,
	}
	_, tokenString, _ := app.TokenAuth.Encode(data) // ใช้ `TokenAuth` สำหรับ JWT
	return tokenString
}

// @Summary User Login
// @Description Handles user login requests and generates a token for the authenticated user.
// @ID usertest-login
// @Tags LoginTest
// @Accept json
// @Produce json
// @Param login-request body request.Login true "User login credentials in JSON format"
// @Success 200 {object} Response{result=string} "JWT token"
// @Failure 400 {object} Response "Bad Request"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /test/login [post]
func (app *Application) LoginTest(w http.ResponseWriter, r *http.Request) {
	// ตรวจสอบ Header
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Invalid content type", http.StatusBadRequest)
		return
	}

	// Decode JSON Payload
	var loginReq req.Login
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received Login Request: %+v\n", loginReq)

	// เรียก service เพื่อเข้าสู่ระบบ
	user, err := app.Service.UserTest.LoginTest(loginReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// สร้าง JWT Token (ถ้าต้องการ)
	token := app.GenerateTokenLogin(user)

	// ส่งกลับข้อมูล
	response := map[string]interface{}{
		"success": true,
		"token":   token,
		"user":    user,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

