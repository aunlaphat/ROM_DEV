package api

import (
	"boilerplate-back-go-2411/logs"
	"boilerplate-back-go-2411/service"

	"github.com/go-chi/jwtauth"
)

type Application struct {
	Logger    logs.Logger
	TokenAuth *jwtauth.JWTAuth
	Service   service.AllOfService
}
