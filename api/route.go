package api

import (
	"net/http"
	"os"
	"path/filepath"

	_ "boilerplate-backend-go/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupSwagger(router *gin.RouterGroup) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func Routes(router *gin.Engine, app *Application) {
	allowedOrigin := os.Getenv("CORS_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:3000"
	}
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{allowedOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Authorization", "Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	apiRouter := router.Group("/api")

	SetupSwagger(apiRouter)

	workDir, err := os.Getwd()
	if err != nil {
		panic("failed to get working directory: " + err.Error())
	}
	filesDir := filepath.Join(workDir, "uploads")
	apiRouter.StaticFS("/uploads", http.Dir(filesDir))

	app.AuthRoute(apiRouter)
	app.UserRoute(apiRouter)
	app.OrderRoute(apiRouter)
	app.DraftConfirmRoute(apiRouter)
}
