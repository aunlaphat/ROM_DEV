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

// 📌 SetupRoutes กำหนดเส้นทาง API ทั้งหมด
func SetupRoutes(router *gin.Engine, app *Application) {
	// Logger and Recovery middleware are already added in server.go
	// router.Use(gin.Logger())
	// router.Use(gin.Recovery())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // อนุญาตเฉพาะ frontend ที่กำหนด
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"X-PINGOTHER", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // 5 นาที
	}))

	apiRouter := router.Group("/api")

	apiRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	/* apiRouter.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to the API!")
	}) */

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "uploads")
	apiRouter.StaticFS("/uploads", http.Dir(filesDir))

	//app.AuthRoute(apiRouter)
	app.OrderRoute(apiRouter)
}
