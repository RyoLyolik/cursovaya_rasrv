package controller

import (
	"httpserver/domain"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AccountController struct {
	Logger         *slog.Logger
	AccountUsecase domain.AccountUsecase
	Token          string
}

func (ac *AccountController) Me(ctx *gin.Context) {
	userData, ok := ctx.Get("user")
	if !ok {
		resp := domain.ErrResp
		resp.Error.Detail = "Unauthorized"
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}
	user := userData.(*domain.User)
	resp := domain.OkResp
	resp.Data = domain.MeOut{
		Email:        user.Email,
		CreationDate: user.CreationDate,
		Role: domain.RoleOut{
			Name: user.Role.Name,
		},
	}
	ctx.JSON(http.StatusOK, resp)
}

func (ac *AccountController) NewAdmin(ctx *gin.Context) {
	var request domain.CreateAdminRequest
	ac.Logger.Info("Incoming request", "ctx", ctx)
	err := ctx.ShouldBind(&request)
	ac.Logger.Info("Binded", "req", request)
	if err != nil {
		resp := domain.ErrResp
		resp.Error.Detail = err.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	if request.Token != ac.Token {
		resp := domain.ErrResp
		resp.Error.Detail = "Invalid token"
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}
	hashedPasswd, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		resp := domain.ErrResp
		resp.Error.Detail = "Failed to hash password"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	user := &domain.User{
		Role:         domain.Admin,
		Email:        request.Email,
		Password:     string(hashedPasswd),
		CreationDate: time.Now(),
	}
	err = ac.AccountUsecase.Create(ctx, user)
	if err != nil {
		resp := domain.ErrResp
		resp.Error.Detail = "Failed to create user"
		ctx.JSON(http.StatusConflict, resp)
		return
	}
	ctx.JSON(http.StatusCreated, domain.OkResp)
}

func (ac *AccountController) NewEmployee(ctx *gin.Context) {
	var request domain.CreateEmployeeRequest
	err := ctx.ShouldBind(&request)
	ac.Logger.Info("Binded", "req", request)
	if err != nil {
		resp := domain.ErrResp
		resp.Error.Detail = "Bad request"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	hashedPasswd, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		resp := domain.ErrResp
		resp.Error.Detail = "Failed to hash password"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	user := &domain.User{
		Role:         domain.Employee,
		Email:        request.Email,
		Password:     string(hashedPasswd),
		CreationDate: time.Now(),
	}
	err = ac.AccountUsecase.Create(ctx, user)
	if err != nil {
		resp := domain.ErrResp
		resp.Error.Detail = "Failed to create user"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	ctx.JSON(http.StatusCreated, domain.OkResp)
}

func (ac *AccountController) List(ctx *gin.Context) {
	users, err := ac.AccountUsecase.List(ctx)
	if err != nil {
		resp := domain.ErrResp
		resp.Error.Detail = "Failed to load users"
		ctx.JSON(http.StatusInternalServerError, resp)
	}
	var resp domain.ListUsersResponse
	resp.Status = "Ok"
	for _, user := range users {
		var respUser domain.UserOut
		respUser.ID = user.ID
		respUser.Email = user.Email
		respUser.CreationDate = user.CreationDate
		respUser.Role = domain.RoleOut{
			Name: user.Role.Name,
		}
		resp.Data = append(resp.Data, respUser)
	}
	ctx.JSON(http.StatusOK, resp)
}
