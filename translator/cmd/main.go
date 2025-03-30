package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"translator/bootstrap"
	messagehandler "translator/message_handler"

	"translator/internal/closer"
)

func main() {
	envFile := flag.String("env-file", "", "Path to environment file")
	flag.Parse()
	app, err := bootstrap.App(envFile)
	if err != nil {
		fmt.Printf("error configuring app %v", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	runApp(ctx, &app)
}

func runApp(ctx context.Context, app *bootstrap.Application) {
	c := &closer.Closer{}
	log := app.Log
	c.Add(app.Shutdown)

	go func() {
		mh := messagehandler.NewMessageHandler(log.WithGroup("MessageHandler"), app.Postgre)
		for {
			_, message, err := app.WS.ReadMessage()
			if err != nil {
				log.Error("read error, trying to reconnect", "err", err)
				if err := app.WS.Close(); err != nil {
					log.Error("failed to close connection", "err", err)
				}
				app.WS, err = bootstrap.Connect(ctx, app.Config.WS)
				if err != nil {
					log.Error("failed to reconnect", "err", err)
				}
				continue
			}
			err = mh.HandleMessage(ctx, message)
			if err != nil {
				log.Error("failed to handle message", "err", err)
			}
		}
	}()

	<-ctx.Done()
	log.Info("Shutting down application")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := c.Close(shutdownCtx); err != nil {
		log.Error("error while shutting down application", "err", err)
	}
}
