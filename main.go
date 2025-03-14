package main

import (
	"boilerplate-back-go-2411/api"
	"boilerplate-back-go-2411/logs"
	"boilerplate-back-go-2411/service"
	"boilerplate-back-go-2411/utils"
	"flag"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-chi/jwtauth"
	"github.com/joho/godotenv"
)

// @title Return Order Management Service ⭐
// @version 1.0
// @description This is a Return Order Management Service API server.
// @BasePath /api

func main() {
	if err := godotenv.Load("./.env"); err != nil {
		log.Fatalln("failed to load .env file: " + err.Error())
	}

	serviceName := "RETURN ORDER MANAGEMENT SERVICE"
	logger, logClose, err := logs.NewLogger(serviceName, "./uploads/error/error.log", 1, 1, 7)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %s", err.Error())
	}
	defer logClose()

	logger.Info(serviceName + " is starting... 🪂 ")

	utils.LoadConfig()

	port := flag.Int("port", utils.AppConfig.ServerPort, "Application server port")
	flag.Parse()
	utils.AppConfig.ServerPort = *port

	SqlDB := utils.GetSqlDB(utils.AppConfig, utils.AppConfig.DatabaseName, *logger)
	defer SqlDB.Close()

	srv := service.NewService(SqlDB, *logger)

	app := &api.Application{
		Logger:    *logger,
		TokenAuth: jwtauth.New("HS256", []byte(utils.AppConfig.JWTSecret), nil),
		Service:   srv,
	}

	if err := app.Serve(); err != nil {
		logger.Fatal(fmt.Sprintf("Server failed to start: %s", err.Error()))
	}
}
