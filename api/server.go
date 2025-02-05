package api

import (
	"boilerplate-backend-go/utils"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// 📌 Serve ทำหน้าที่ Start Server และ Handle Graceful Shutdown
func (app *Application) Serve() error {
	// ✅ สร้าง Router (Gin Engine) โดยไม่ใช้ gin.Default()
	router := gin.New()

	// ✅ เพิ่ม Logger และ Recovery Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// ✅ Setup Routes
	SetupRoutes(router, app)

	// ✅ กำหนดพอร์ต
	serverPort := fmt.Sprintf(":%d", utils.AppConfig.ServerPort)

	// ✅ สร้าง HTTP Server
	srv := &http.Server{
		Addr:         serverPort,
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// ✅ Graceful Shutdown
	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.Logger.Info(fmt.Sprintf("🛑 Shutting down server with signal: %s", s.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		shutdownError <- srv.Shutdown(ctx)
	}()

	app.Logger.Info(fmt.Sprintf("🚀 Starting server at port: %d", utils.AppConfig.ServerPort))

	// ✅ เริ่มให้เซิร์ฟเวอร์รับ Request
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// ✅ รอให้ Shutdown เสร็จสิ้น
	err = <-shutdownError
	if err != nil {
		return err
	}

	app.Logger.Info(fmt.Sprintf("🛑 Server stopped at port: %d", utils.AppConfig.ServerPort))

	return nil
}
