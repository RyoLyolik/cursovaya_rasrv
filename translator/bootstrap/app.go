package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

type Application struct {
	Config   *Config
	Log      *slog.Logger
	Postgre  *sql.DB
	WSClient *websocket.Conn
}

func App(envFile *string) (Application, error) {
	app := &Application{}
	app.Config = MustLoad(envFile)
	app.Log = SetupLogger(&app.Config.LogCfg)

	app.Log.Info("Starting application")
	app.Log.Debug("Debug messages are enabled")

	psqlConnectionInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		app.Config.SQLStorage.Host,
		app.Config.SQLStorage.Port,
		app.Config.SQLStorage.User,
		app.Config.SQLStorage.Password,
		app.Config.SQLStorage.DBName,
	)
	postgre, err := MakePostgres(app.Log, psqlConnectionInfo)
	if err != nil {
		return *app, err
	}
	app.Postgre = postgre
	app.WSClient, err = WSConnect(context.Background(), app.Config.WS)
	if err != nil {
		return *app, err
	}
	return *app, nil
}

func (a *Application) WSReconnect(ctx context.Context) error {
	if a.WSClient != nil {
		err := a.WSClient.Close()
		if err != nil {
			return err
		}
	}
	var err error
	a.WSClient, err = WSConnect(ctx, a.Config.WS)
	if err != nil {
		return err
	}
	return nil
}

func (a *Application) Shutdown(ctx context.Context) error {
	a.Log.Info("Shutting down application")

	a.Log.Info("Closing postgres connection")
	if err := a.Postgre.Close(); err != nil {
		return err
	}
	a.Log.Info("Closing websocket connection")
	if err := a.WSClient.Close(); err != nil {
		return err
	}
	a.Log.Info("Application stopped successfully")
	return nil
}
