package main

import (
	"context"
	"flag"
	"fmt"
	"httpserver/api/route"
	"httpserver/bootstrap"
	"httpserver/internal/closer"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	envFile := flag.String("env-file", "", "Path to environment file")
	flag.Parse()
	if *envFile == "" {
		fmt.Println("error: --env-file argument is required")
		flag.Usage()
		os.Exit(1)
	}
	app, err := bootstrap.App(envFile)
	if err != nil {
		fmt.Printf("error configuring app%v", err)
		os.Exit(1)
	}

	// route.Setup()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	runApp(ctx, &app)
}

func runApp(ctx context.Context, app *bootstrap.Application) {
	c := &closer.Closer{}
	log := app.Log
	gin := gin.Default()
	address := fmt.Sprintf("%s:%d", app.Config.HttpServer.Host, app.Config.HttpServer.Port)

	c.Add(app.Shutdown)

	route.Setup(app, gin)

	log.Info("Starting server")
	go func() {
		gin.Run(address)
	}()
	log.Info("Listening", "address", address)

	<-ctx.Done()
	log.Info("Shutting down application")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := c.Close(shutdownCtx); err != nil {
		log.Error("error while shutting down application", "err", err)
	}
}
