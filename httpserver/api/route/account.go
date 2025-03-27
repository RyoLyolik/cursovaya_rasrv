package route

import (
	"database/sql"
	"httpserver/api/controller"
	"httpserver/repository"
	"httpserver/usecase"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func NewProtectedAccountRouter(log *slog.Logger, group *gin.RouterGroup, db *sql.DB, sessionStorage *redis.Client, timeout time.Duration) {
	ur := repository.NewUserRepository(db, log.With("repository", "user"))
	sr := repository.NewSessionRepository(sessionStorage, log.With("repository", "session"), timeout)
	ac := &controller.AccountController{
		Logger:         log.With("controller", "login"),
		AccountUsecase: usecase.NewAccountUsecase(ur, sr, timeout),
	}
	group.GET("/account/me", ac.Me)
}

func NewPublicAccountRouter(log *slog.Logger, group *gin.RouterGroup, db *sql.DB, timeout time.Duration, token string) {
	ur := repository.NewUserRepository(db, log.With("repository", "user"))
	ac := &controller.AccountController{
		Logger:         log.With("controller", "login"),
		AccountUsecase: usecase.NewAccountUsecase(ur, nil, timeout),
		Token:          token,
	}
	group.POST("/admin", ac.NewAdmin)
}

func NewAdminRouter(log *slog.Logger, group *gin.RouterGroup, db *sql.DB, sessionStorage *redis.Client, timeout time.Duration) {
	ur := repository.NewUserRepository(db, log.With("repository", "user"))
	sr := repository.NewSessionRepository(sessionStorage, log.With("repository", "session"), timeout)
	ac := &controller.AccountController{
		Logger:         log.With("controller", "login"),
		AccountUsecase: usecase.NewAccountUsecase(ur, sr, timeout),
	}
	group.POST("/account", ac.NewEmployee)
}
