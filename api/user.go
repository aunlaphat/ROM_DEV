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

/*
// AddUser godoc
// @Summary Add a new user
// @Description Add a new user to the system
// @ID add-user
// @Tags User
// @Accept json
// @Produce json
// @Param user body request.User true "User details"
// @Success 201 {object} api.Response "User added successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /user/add-user [post]
func (app *Application) AddUser(w http.ResponseWriter, r *http.Request) {
	var req request.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := app.Service.User.AddUser(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "User added successfully", nil, http.StatusCreated)
}

// EditUser godoc
// @Summary Edit an existing user
// @Description Edit the role and warehouse of an existing user
// @ID edit-user
// @Tags User
// @Accept json
// @Produce json
// @Param user body request.User true "User details"
// @Success 200 {object} api.Response "User edited successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /user/edit-user [put]
func (app *Application) EditUser(w http.ResponseWriter, r *http.Request) {
	var req request.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := app.Service.User.EditUser(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "User edited successfully", nil, http.StatusOK)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user from the system
// @ID delete-user
// @Tags User
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {object} api.Response "User deleted successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /user/delete-user/{userID} [delete]
func (app *Application) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		http.Error(w, "UserID is required", http.StatusBadRequest)
		return
	}

	err := app.Service.User.DeleteUser(r.Context(), userID)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "User deleted successfully", nil, http.StatusOK)
}
*/
