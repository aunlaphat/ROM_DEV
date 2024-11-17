package api

import (
	"boilerplate-backend-go/logs"
	"boilerplate-backend-go/service"

	"github.com/go-chi/jwtauth"
)

type Config struct {
	Port        int
	SmsApiKey   string
	SmsClientId string
	SmsSenderId string
}

type Application struct {
	Config    Config
	Logger    logs.Logger
	TokenAuth *jwtauth.JWTAuth
	Service   service.AllOfService
}
