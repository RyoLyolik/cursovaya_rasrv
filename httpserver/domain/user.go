package domain

import (
	"context"
	"time"
)

type User struct {
	ID           int
	Email        string
	Password     string
	CreationDate time.Time
	Role         Role
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id int) (*User, error)
	List(ctx context.Context) ([]*User, error)
}
