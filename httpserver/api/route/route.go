package route

import (
	"httpserver/api/middleware"
	"httpserver/bootstrap"
	"time"

	"github.com/gin-gonic/gin"
)

func Setup(app *bootstrap.Application, gin *gin.Engine) {
	log := app.Log.WithGroup("http-server")
	log.Info("Setup routes")

	mw := middleware.NewMiddleWare(
		log.With("middleware", "session_auth"),
		app.Postgre,
		app.SessionStorage,
		5*time.Second,
	)

	publicRouter := gin.Group("")
	NewLoginRouter(
		log.With("route", "public"),
		publicRouter,
		app.Postgre,
		app.SessionStorage,
		5*time.Second,
	)
	NewPublicAccountRouter(
		log.With("route", "public"),
		publicRouter,
		app.Postgre,
		5*time.Second,
		app.Config.AdminToken,
	)

	protectedRouter := gin.Group("")
	protectedRouter.Use(mw.AuthMiddleware())
	NewProtectedAccountRouter(
		log.With("route", "proteced"),
		protectedRouter,
		app.Postgre,
		app.SessionStorage,
		5*time.Second,
	)

	adminRouter := protectedRouter.Group("")
	adminRouter.Use(mw.CheckAdminMiddleware())
	NewAdminRouter(
		log.With("route", "proteced"),
		adminRouter,
		app.Postgre,
		app.SessionStorage,
		5*time.Second,
	)
}
