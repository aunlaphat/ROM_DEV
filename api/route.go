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
		allowedOrigin = "http://localhost:3000" // Default
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{allowedOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"X-PINGOTHER", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // 5 minutes
	}))

	apiRouter := router.Group("/api")

	SetupSwagger(apiRouter)

	workDir, err := os.Getwd()
	if err != nil {
		panic("failed to get working directory: " + err.Error())
	}
	filesDir := filepath.Join(workDir, "uploads")
	apiRouter.StaticFS("/uploads", http.Dir(filesDir))

	// Authenticated & User Routes
	app.AuthRoute(apiRouter)
	app.UserRoute(apiRouter)

	app.OrderRoute(apiRouter)
	app.DraftConfirmRoute(apiRouter)
}
