package controller

import (
	"httpserver/domain"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginController struct {
	Logger       *slog.Logger
	LoginUsecase domain.LoginUsecase
	Expiration   int
}

func (lc *LoginController) Login(ctx *gin.Context) {
	var request domain.LoginRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		resp := domain.ErrResp
		resp.Error.Detail = err.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	user, err := lc.LoginUsecase.GetUserByEmail(ctx, request.Email)
	if err != nil {
		lc.Logger.Info("Accont was not found", "err", err)
		resp := domain.ErrResp
		resp.Error.Detail = "User not found"
		ctx.JSON(http.StatusNotFound, resp)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)) != nil {
		resp := domain.ErrResp
		resp.Error.Detail = "Invalid credentials"
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}
	sessionID, err := lc.LoginUsecase.SaveSession(ctx, user.ID)
	if err != nil {
		resp := domain.ErrResp
		resp.Error.Detail = err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	ctx.SetCookie("brick-session-id", sessionID, lc.Expiration, "*", "*", true, true)
	ctx.JSON(http.StatusOK, domain.OkResp)
}
