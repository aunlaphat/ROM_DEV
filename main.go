package main

import (
	"boilerplate-backend-go/api"
	"boilerplate-backend-go/logs"
	"boilerplate-backend-go/service"
	"boilerplate-backend-go/utils"
	"flag"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-chi/jwtauth"
	"github.com/joho/godotenv"
)

// @title Return Order Management Service API ⭐
// @version 1.0
// @description This is a Return Order Management Service API server.
// @BasePath /api

func main() {
	// ✅ โหลด Environment Variables
	if err := godotenv.Load("./.env"); err != nil {
		log.Fatalln("failed to load .env file: " + err.Error())
	}

	// ✅ กำหนดค่าพารามิเตอร์ของ Server
	serviceName := "RETURN ORDER MANAGEMENT SERVICE"
	logger, logClose, err := logs.NewLogger(serviceName, "./uploads/error/error.log", 1, 1, 7)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %s", err.Error())
	}
	defer logClose() // ✅ ปิด Logger เมื่อโปรแกรมสิ้นสุด

	logger.Info(serviceName + " is starting... 🪂 ")

	// ✅ โหลดค่าคอนฟิกจาก `.env`
	utils.LoadConfig()

	// ✅ ใช้ `flag` กำหนดค่าพอร์ตจาก CLI
	port := flag.Int("port", utils.AppConfig.ServerPort, "Application server port")
	flag.Parse()
	utils.AppConfig.ServerPort = *port // ✅ ใช้ค่าพอร์ตจาก Flag

	// ✅ เชื่อมต่อกับฐานข้อมูล SQL Server
	SqlDB := utils.GetSqlDB(utils.AppConfig, utils.AppConfig.DatabaseName, *logger)
	defer SqlDB.Close()

	// ✅ สร้าง Service Layer
	srv := service.NewService(SqlDB, *logger)

	// ✅ สร้าง Application Instance
	app := &api.Application{
		Logger:    *logger,
		TokenAuth: jwtauth.New("HS256", []byte(utils.AppConfig.JWTSecret), nil),
		Service:   srv,
	}

	// ✅ เริ่มต้นเซิร์ฟเวอร์
	if err := app.Serve(); err != nil {
		logger.Fatal(fmt.Sprintf("Server failed to start: %s", err.Error()))
	}
}
