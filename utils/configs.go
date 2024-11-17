package utils

// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"strconv"

// 	"github.com/joho/godotenv"
// )

// type ConfigMongo struct {
// 	Url    string
// 	NameDB string
// }

// var AppConfig Config

// func LoadConfigMongo() (*ConfigMongo, error) {
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 		return nil, err
// 	}

// 	subfix := ""
// 	if os.Getenv("MODE") == "DEV" {
// 		subfix = "_DEV"
// 		fmt.Println("DEV")
// 	}

// 	configMongo := &ConfigMongo{
// 		Url:    os.Getenv("MONGOURI" + subfix),
// 		NameDB: os.Getenv("MONGO_DB_NAME" + subfix),
// 	}

// 	return configMongo, nil
// }

// type ConfigSQL struct {
// 	DBHost     string
// 	DBPort     string
// 	DataBase   string
// 	DBUser     string
// 	DBPassword string
// }

// func LoadConfigSQL() (*ConfigSQL, error) {
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 		return nil, err
// 	}
// 	subfix := ""
// 	if os.Getenv("MODE") == "DEV" {
// 		subfix = "_DEV"
// 	}

// 	// Create connection string
// 	configSQL := &ConfigSQL{
// 		DBHost:     os.Getenv("MSSQL_DB_SERVER" + subfix),
// 		DBPort:     os.Getenv("MSSQL_DB_PORT" + subfix),
// 		DataBase:   os.Getenv("MSSQL_DB_DATABASE" + subfix),
// 		DBUser:     os.Getenv("MSSQL_DB_USER" + subfix),
// 		DBPassword: os.Getenv("MSSQL_DB_PASSWORD" + subfix),
// 	}
// 	return configSQL, nil
// }

// func LoadPort() (int, error) {
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		log.Fatal("Error loading .env file 1")
// 		return 0, err
// 	}
// 	subfix := ""
// 	if os.Getenv("MODE") == "DEV" {
// 		subfix = "_DEV"
// 	}
// 	port := os.Getenv("PORT" + subfix)

// 	port_num, err := strconv.Atoi(port)
// 	if err != nil {
// 		log.Fatal("Error loading .env file 2")
// 		return 0, err
// 	}
// 	return port_num, nil
// }

// func LoadAxApi() (string, error) {
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		log.Fatal("Error loading .env file 3")
// 		return "", err
// 	}
// 	subfix := ""
// 	if os.Getenv("MODE") == "DEV" {
// 		subfix = "_DEV"
// 	}
// 	ax_api := os.Getenv("AX_API" + subfix)
// 	return ax_api, nil
// }

// func LoadJWTSerect() (string, error) {
// 	err := godotenv.Load(".env")
// 	if err != nil {

// 		log.Fatal("Error loading .env file 4")
// 		return "", err
// 	}
// 	subfix := ""
// 	if os.Getenv("MODE") == "DEV" {
// 		subfix = "_DEV"
// 	}
// 	ax_api := os.Getenv("JWT_SECRET" + subfix)
// 	return ax_api, nil
// }

// type ConfigSMS struct {
// 	SMSApiKey   string
// 	SMSClientID string
// 	SMSSenderID string
// }

// func checkEnv() error {
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		log.Fatal("Error loading .env file 5")
// 		return err
// 	}
// 	return nil

// }
// func LoadConfigSMS() (*ConfigSMS, error) {
// 	err := checkEnv()
// 	if err != nil {
// 		return nil, err
// 	}

// 	configSMS := &ConfigSMS{
// 		SMSApiKey:   os.Getenv("SMS_APIKEY"),
// 		SMSClientID: os.Getenv("SMS_CLIENT_ID"),
// 		SMSSenderID: os.Getenv("SMS_SENDER_ID"),
// 	}
// 	return configSMS, nil
// }

// func LoadConfigGoogleMapAPI() (string, error) {
// 	err := checkEnv()
// 	if err != nil {
// 		return "", err
// 	}
// 	subfix := ""
// 	if os.Getenv("MODE") == "DEV" {
// 		subfix = "_DEV"
// 	}
// 	googleMapApi := os.Getenv("GOOGLE_MAP_APIKEY" + subfix)
// 	return googleMapApi, nil
// }

// type ConfigRedis struct {
// 	Address  string
// 	Password string
// 	DB       int
// }

// func LoadConfigRedis() (*ConfigRedis, error) {

// 	err := checkEnv()
// 	if err != nil {
// 		return nil, err
// 	}
// 	db, _ := strconv.Atoi(os.Getenv("REDIS_DB_DATABASE"))
// 	configRedis := &ConfigRedis{
// 		Address:  os.Getenv("REDIS_DB_ADDRESS"),
// 		Password: os.Getenv("REDIS_DB_PASSWORD"),
// 		DB:       db,
// 	}
// 	return configRedis, nil
// }
