package domain

import (
	"context"
	"time"
)

type MeOut struct {
	Email        string    `json:"email"`
	CreationDate time.Time `json:"creation_date" binding:"datetime"`
	Role         RoleOut   `json:"role"`
}

type MeResponse struct {
	DefaultResponse
	Data MeOut `json:"data,omitempty"`
}

type CreateEmployeeRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateAdminRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Token    string `json:"token" binding:"required"`
}

type UserOut struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	CreationDate time.Time `json:"creation_date" binding:"datetime"`
	Role         RoleOut   `json:"role"`
}

type ListUsersResponse struct {
	DefaultResponse
	Data []UserOut `json:"data"`
}

type AccountUsecase interface {
	AccountBySession(ctx context.Context, session string) (*User, error)
	Create(ctx context.Context, user *User) error
	List(ctx context.Context) ([]*User, error)
}
