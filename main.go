package main

import (

	// "log"

	"os"
	"os/signal"
	"time"

	"github.com/rewiko/boilerplate-echo/server"
)

func main() {
	e := server.Start()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	// healthCheck = "unhealthy"
	time.Sleep(1 * time.Second)
	server.Stop(e)
}
