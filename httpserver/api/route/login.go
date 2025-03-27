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

func NewLoginRouter(log *slog.Logger, group *gin.RouterGroup, db *sql.DB, sessionStorage *redis.Client, timeout time.Duration) {
	tokenEx := time.Hour * 24 * 30 // TODO move in config
	ur := repository.NewUserRepository(db, log.With("repository", "user"))
	sr := repository.NewSessionRepository(sessionStorage, log.With("repository", "session"), tokenEx)
	lc := &controller.LoginController{
		Logger:       log.With("controller", "login"),
		LoginUsecase: usecase.NewLoginUsecase(ur, sr, timeout),
		Expiration:   int(tokenEx) / int(time.Second),
	}
	group.POST("/login", lc.Login)
}
