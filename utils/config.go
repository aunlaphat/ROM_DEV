package utils

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	SQLHost          string
	SQLPort          string
	DatabaseName     string
	SQLUser          string
	SQLPassword      string
	ServerPort       int
	AXApi            string
	JWTSecret        string
	SMSApiKey        string
	SMSClientID      string
	SMSSenderID      string
	LarkMSGType      string
	LarkAppID        string
	LarkAppSecret    string
	LarkApprovalCode string
}

var AppConfig Config

func LoadConfig() {

	AppConfig = Config{
		SQLHost:          getEnv("MSSQL_DB_SERVER", ""),
		SQLPort:          getEnv("MSSQL_DB_PORT", ""),
		DatabaseName:     getEnv("MSSQL_DB_DATABASE", ""),
		SQLUser:          getEnv("MSSQL_DB_USER", ""),
		SQLPassword:      getEnv("MSSQL_DB_PASSWORD", ""),
		ServerPort:       getEnvAsInt("PORT", 8080),
		AXApi:            getEnv("AX_API", ""),
		JWTSecret:        getEnv("JWT_SECRET", ""),
		SMSApiKey:        getEnv("SMS_APIKEY", ""),
		SMSClientID:      getEnv("SMS_CLIENT_ID", ""),
		SMSSenderID:      getEnv("SMS_SENDER_ID", ""),
		LarkMSGType:      getEnv("LARK_MSG_TYPE", ""),
		LarkAppID:        getEnv("LARK_APP_ID", ""),
		LarkAppSecret:    getEnv("LARK_APP_SECRET", ""),
		LarkApprovalCode: getEnv("LARK_APPROVE_CODE", ""),
	}
}

var keyExist = "Key exist "

// getEnv gets the environment variable key or returns the defaultValue if not found
func getEnv(key string, defaultValue string) string {

	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Fatalln(keyExist + key)
	return defaultValue
}

// getEnvAsInt gets the environment variable as an integer or returns the defaultValue
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	log.Fatalln(keyExist + key)

	return defaultValue
}

// // getEnvAsBool gets the environment variable as a boolean or returns the defaultValue
// func getEnvAsBool(key string, defaultValue bool) bool {
// 	valueStr := getEnv(key, "")
// 	if value, err := strconv.ParseBool(valueStr); err == nil {
// 		return value
// 	}
// 	log.Fatalln(keyExist + key)

// 	return defaultValue
// }
