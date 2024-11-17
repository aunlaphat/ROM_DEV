package api

import (
	"boilerplate-backend-go/logs"
	"boilerplate-backend-go/service"

	"github.com/go-chi/jwtauth"
)

type Application struct {
	Logger    logs.Logger
	TokenAuth *jwtauth.JWTAuth
	Service   service.AllOfService
}
