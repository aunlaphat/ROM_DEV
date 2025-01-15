package main

import (
	api "boilerplate-backend-go/api"
	"boilerplate-backend-go/logs"
	"boilerplate-backend-go/service"
	"boilerplate-backend-go/utils"
	"flag"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-chi/jwtauth"
	"github.com/joho/godotenv"
)

// @title Boilerplate Service
// @version 1.0
// @description This is a sample server for Boilerplate project .
// contact.name API Support

// @BasePath /api
func main() {
	//Env loading
	err := godotenv.Load("./.env")

	if err != nil {
		log.Fatalln("failed to load .env file: " + err.Error())
	}

	//Flag variable
	serviceName := "RETURN ORDER SERVICE"
	logger, logClose, err := logs.NewLogger(serviceName, "./uploads/error/error.log", 1, 1, 7)
	defer logClose()
	if err != nil {
		panic(err)
	}

	//Instance logger of service
	logger.Info(serviceName + " is starting...")

	utils.LoadConfig()

	flag.IntVar(&utils.AppConfig.ServerPort, "port", utils.AppConfig.ServerPort, "Application server port")
	flag.Parse()

	//SQL Server database
	SqlDB := utils.GetSqlDB(utils.AppConfig, utils.AppConfig.DatabaseName, *logger)

	defer SqlDB.Close()
	//SMS

	//Instance Service
	srv := service.NewService(SqlDB, *logger)
	//Instance application
	app := &api.Application{
		Logger:    *logger,
		TokenAuth: jwtauth.New("HS256", []byte(utils.AppConfig.JWTSecret), nil),
		Service:   srv,
	}
	// Run application
	err = app.Serve()
	if err != nil {
		logger.Fatal(err.Error())
	}

}
