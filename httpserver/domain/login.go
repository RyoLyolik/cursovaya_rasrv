package domain

import (
	"context"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	DefaultResponse
	Data struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	} `json:"data,omitempty"`
}

type LoginUsecase interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	SaveSession(ctx context.Context, id int) (string, error)
}
