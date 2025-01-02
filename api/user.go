package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// UserRoute defines the routes for user operations
func (app *Application) UserRoute(apiRouter *chi.Mux) {
	apiRouter.Route("/user", func(r chi.Router) {
		r.Get("/get-user", app.GetUser) // New route for getting user
		r.Get("/get-user-with-permission", app.GetUserWithPermission)
	})
}

// GetUser godoc
// @Summary Get user by username and password
// @Description Retrieve the details of a user by their username and password
// @ID get-user
// @Tags User
// @Accept json
// @Produce json
// @Param username query string true "Username"
// @Param password query string true "Password"
// @Success 200 {object} api.Response{data=response.Login} "User retrieved successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 404 {object} api.Response "User not found"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /user/get-user [get]
func (app *Application) GetUser(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	user, err := app.Service.User.GetUser(r.Context(), username, password)
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
// @Param username query string true "Username"
// @Param password query string true "Password"
// @Success 200 {object} api.Response{data=response.UserPermission} "User with permissions retrieved successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 404 {object} api.Response "User not found"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /user/get-user-with-permission [get]
func (app *Application) GetUserWithPermission(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	user, err := app.Service.User.GetUserWithPermission(r.Context(), username, password)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "User with permissions retrieved successfully", user, http.StatusOK)
}
