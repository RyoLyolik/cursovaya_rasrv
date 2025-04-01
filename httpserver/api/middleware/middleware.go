package middleware

import (
	"database/sql"
	"httpserver/domain"
	"httpserver/repository"
	"httpserver/usecase"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type MiddleWare struct {
	Logger         *slog.Logger
	AccountUsecase domain.AccountUsecase
}

func NewMiddleWare(log *slog.Logger, db *sql.DB, sessionStorage *redis.Client, timeout time.Duration) *MiddleWare {
	ur := repository.NewUserRepository(db, log.With("repository", "user"))
	sr := repository.NewSessionRepository(sessionStorage, log.With("repository", "session"), timeout)
	samc := &MiddleWare{
		Logger:         log.With("controller", "login"),
		AccountUsecase: usecase.NewAccountUsecase(ur, sr, timeout),
	}
	return samc
}

func (samc *MiddleWare) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Request.Cookie("brick-session-id")
		if err != nil {
			resp := domain.ErrResp
			resp.Error.Detail = "Cannot validate session"
			ctx.JSON(http.StatusBadRequest, resp)
			ctx.Abort()
			return
		}
		session := cookie.Value
		user, err := samc.AccountUsecase.AccountBySession(ctx, session)
		if err != nil {
			resp := domain.ErrResp
			resp.Error.Detail = "Unauthorized"
			ctx.JSON(http.StatusUnauthorized, resp)
			ctx.Abort()
			return
		}
		ctx.Set("user", user)
		ctx.Next()
	}
}

func (samc *MiddleWare) CheckAdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userData, ok := ctx.Get("user")
		if !ok {
			resp := domain.ErrResp
			resp.Error.Detail = "Unauthorized"
			ctx.JSON(http.StatusUnauthorized, resp)
			return
		}
		user := userData.(*domain.User)
		if user.Role != domain.Admin {
			resp := domain.ErrResp
			resp.Error.Detail = "Not permitted"
			ctx.JSON(http.StatusForbidden, resp)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
