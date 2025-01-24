package api

import (
	"boilerplate-backend-go/dto/request"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// UserRoute defines the routes for user operations
func (app *Application) UserRoute(apiRouter *chi.Mux) {
	apiRouter.Route("/user", func(r chi.Router) {
		r.Post("/get-user", app.GetUser)
		r.Post("/get-user-with-permission", app.GetUserWithPermission)
		/* 		r.Post("/add-user", app.AddUser)
		   		r.Put("/edit-user", app.EditUser)
		   		r.Delete("/delete-user/{userID}", app.DeleteUser) */
	})
}

// GetUser godoc
// @Summary Get user by username and password
// @Description Retrieve the details of a user by their username and password
// @ID get-user
// @Tags User
// @Accept json
// @Produce json
// @Param LoginWeb body request.LoginWeb true "User login credentials in JSON format"
// @Success 200 {object} api.Response{data=response.Login} "User retrieved successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 404 {object} api.Response "User not found"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /user/get-user [post]
func (app *Application) GetUser(w http.ResponseWriter, r *http.Request) {
	var req request.LoginWeb
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.UserName == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	user, err := app.Service.User.GetUser(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "User retrieved successfully", user, http.StatusOK)
}

// GetUserWithPermission godoc
// @Summary Get user with permissions by username and password
// @Description Retrieve the details of a user with permissions by their username and password
// @ID get-user-with-permission
// @Tags User
// @Accept json
// @Produce json
// @Param login body request.LoginLark true "User login credentials in JSON format"
// @Success 200 {object} api.Response{data=response.UserPermission} "User with permissions retrieved successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 404 {object} api.Response "User not found"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /user/get-user-with-permission [post]
func (app *Application) GetUserWithPermission(w http.ResponseWriter, r *http.Request) {
	var req request.LoginLark
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.UserID == "" || req.UserName == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	user, err := app.Service.User.GetUserWithPermission(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "User with permissions retrieved successfully", user, http.StatusOK)
}
