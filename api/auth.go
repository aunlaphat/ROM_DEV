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

func (app *Application) AuthRoute(apiRouter *chi.Mux) {

	apiRouter.Route("/auth", func(r chi.Router) {
		r.Post("/login", app.Login)
		r.Post("/login-lark", app.LoginFromLark)
		r.Group(func(router chi.Router) {
			router.Use(jwtauth.Verifier(app.TokenAuth))
			router.Use(jwtauth.Authenticator)
			router.Get("/", app.CheckAuthen)
			router.Post("/logout", app.Logout)
		})
	})
}

var contentType = "content-type"
var appJson = "application/json"

// Generate JWT token with username's payload
func (app *Application) GenerateToken(tokenData res.Login) string {
	data := map[string]interface{}{
		"userID":     tokenData.UserID,
		"roleID":     tokenData.RoleID,
		"nickName":   tokenData.NickName,
		"fullNameTH": tokenData.FullNameTH,
		"plateform":  tokenData.Platform,
	}
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
func (app *Application) Login(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(contentType) != appJson {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	req := req.LoginWeb{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(w, err)
		return
	}

	user, err := app.Service.User.Login(req)
	if err != nil {
		handleError(w, err)
		return
	}
	tokenData := res.Login{
		UserID:     user.UserID,
		RoleID:     user.RoleID,
		NickName:   user.NickName,
		FullNameTH: user.FullNameTH,
		Platform:   user.Platform,
	}
	fmt.Println("token data", tokenData)
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
		UserID:     user.UserID,
		RoleID:     user.RoleID,
		NickName:   user.NickName,
		FullNameTH: user.FullNameTH,
		Platform:   user.Platform,
	}
	token := app.GenerateToken(tokenData)

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
func (app *Application) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		MaxAge:   -1, // Delete the cookie.
		SameSite: http.SameSiteLaxMode,
		Name:     "jwt",
		Value:    "",
	})
	handleResponse(w, true, "login Success", "None", http.StatusOK)

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
func (app *Application) CheckAuthen(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Context())
	_, claims, _ := jwtauth.FromContext(r.Context())
	handleResponse(w, true, "Checked", claims, http.StatusOK)
}
