package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type Application struct {
	Config         *Config
	Log            *slog.Logger
	Postgre        *sql.DB
	SessionStorage *redis.Client
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
		return *app, fmt.Errorf("%v", err)
	}
	app.Postgre = postgre

	redisAddr := fmt.Sprintf("%s:%s", app.Config.SessionStorage.Host, app.Config.SessionStorage.Port)
	redis, err := MakeRedis(app.Log, redisAddr)
	if err != nil {
		return *app, fmt.Errorf("%v", err)
	}
	app.SessionStorage = redis

	return *app, nil
}

func (a *Application) Shutdown(ctx context.Context) error {
	a.Log.Info("Shutting down application")

	a.Log.Info("Closing postgres connection")
	if err := a.Postgre.Close(); err != nil {
		return err
	}

	a.Log.Info("Closing redis")
	if err := a.SessionStorage.Close(); err != nil {
		return err
	}

	a.Log.Info("Application stopped successfully")
	return nil
}
