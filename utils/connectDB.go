package utils

import (
	"boilerplate-backend-go/logs"
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMongoDB(Url string, database string, logger logs.Logger) *mongo.Database {
	clientOptions := options.Client().ApplyURI(Url)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Fatal(err.Error())
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Fatal(err.Error())
	}

	mongoDB := client.Database(database)
	logger.Info("Connected to MongoDB ")

	return mongoDB
}

func GetSqlDB(config ConfigSQL, database string, logger logs.Logger) *sqlx.DB {
	connStringDatabase := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;trustservercertificate=true;encrypt=DISABLE",
		config.DBHost, config.DBUser, config.DBPassword, config.DBPort, database)
	dbTemp, err := sqlx.Open("sqlserver", connStringDatabase)
	if err != nil {
		logger.Fatal(err.Error())
	}
	err = dbTemp.Ping()
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info("Connected to MSSQL at " + config.DBHost)

	return dbTemp
}
