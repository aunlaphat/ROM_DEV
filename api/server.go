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

func (app *Application) Serve() error {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	Routes(router, app)

	serverPort := fmt.Sprintf(":%d", utils.AppConfig.ServerPort)

	srv := &http.Server{
		Addr:         serverPort,
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.Logger.Info(fmt.Sprintf("Shutting down server with signal: %s", s.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		shutdownError <- srv.Shutdown(ctx)
	}()

	app.Logger.Info(fmt.Sprintf("Starting server at port: %d", utils.AppConfig.ServerPort))

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.Logger.Info(fmt.Sprintf("Server stopped at port: %d", utils.AppConfig.ServerPort))

	return nil
}
