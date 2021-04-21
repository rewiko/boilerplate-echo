package server

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	customMid "github.com/rewiko/boilerplate-echo/middleware"

	// "log"
	"net/http"
	"time"
)

// Stop - stop the echo server
func Stop(e *echo.Echo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	return nil
}

// Start return echo config and the http server
func Start() *echo.Echo {

	// Echo instance
	e := echo.New()
	// e.Logger.SetLevel(log.ERROR)
	e.Logger.SetLevel(log.INFO)

	// Middleware
	e.Use(customMid.NewMetric())

	logSkipper := func(c echo.Context) bool {
		if c.Response().Status == 200 {
			return true
		}
		return false
	}
	e.Use(customMid.LoggerWithConfig(customMid.LoggerConfig{
		Skipper: logSkipper,
	}))

	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("5M"))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	e.GET("/", hello)

	// Start server
	s := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}

	//s.SetKeepAlivesEnabled(false)
	e.HideBanner = true

	// Start server
	go func() {
		if err := e.StartServer(s); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	return e
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello !")
}
