package main

import (
	api "boilerplate-backend-go/cmd/api"
	"boilerplate-backend-go/logs"
	"boilerplate-backend-go/service"
	"boilerplate-backend-go/utils"
	"flag"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-chi/jwtauth"
	"github.com/joho/godotenv"
)

// @title FOC
// @version 1.0
// @description This is a sample server for FOC project not Frog OOB OOB.
// contact.name API Support

// @BasePath /api
func main() {
	//Flag variable
	logger, logClose, err := logs.NewLogger("./uploads/error/error.log", 1, 1, 7)
	defer logClose()
	if err != nil {
		panic(err)
	}

	//Instance logger of service
	logger.Info("TRADE ORDER SERVICE")

	var cfg api.Config
	port_env, err := utils.LoadPort()
	if err != nil {
		logger.Fatal(err.Error())
		port_env = 8080
	}

	flag.IntVar(&cfg.Port, "port", port_env, "Application server port")
	flag.Parse()

	//Env loading
	err = godotenv.Load("./.env")

	if err != nil {
		logger.Fatal(err.Error())
	}
	//JWT Serect
	jwtSerect, err := utils.LoadJWTSerect()
	if err != nil {
		logger.Fatal(err.Error())
	}

	//SQL Server database
	configSql, err := utils.LoadConfigSQL()
	if err != nil {
		logger.Fatal(err.Error())
	}
	SqlDB := utils.GetSqlDB(*configSql, configSql.DataBase, *logger)

	defer SqlDB.Close()
	//SMS
	configSms, err := utils.LoadConfigSMS()
	if err != nil {
		logger.Fatal(err.Error())
	}
	cfg.SmsApiKey = configSms.SMSApiKey
	cfg.SmsClientId = configSms.SMSClientID
	cfg.SmsSenderId = configSms.SMSSenderID

	//Redis
	// configRedis, err := utils.LoadConfigRedis()
	// if err != nil {
	// 	logger.Fatal(err.Error())
	// }
	// redisDB := redis.NewClient(&redis.Options{
	// 	Addr:     configRedis.Address,
	// 	Password: configRedis.Password,
	// 	DB:       configRedis.DB,
	// })

	//Instance Service
	srv := service.NewService(SqlDB, *logger)
	//Instance application
	app := &api.Application{
		Config:    cfg,
		Logger:    *logger,
		TokenAuth: jwtauth.New("HS256", []byte(jwtSerect), nil),
		Service:   srv,
	}
	// Run application
	err = app.Serve()
	if err != nil {
		logger.Fatal(err.Error())
	}

}
