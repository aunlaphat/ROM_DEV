package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// UserRoute defines the routes for user operations
func (app *Application) UserRoute(apiRouter *chi.Mux) {
	apiRouter.Route("/user", func(r chi.Router) {
		r.Get("/get-user", app.GetUser) // New route for getting user
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
// @Success 200 {object} api.Response
// @Failure 404 {object} api.Response
// @Failure 500 {object} api.Response
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
