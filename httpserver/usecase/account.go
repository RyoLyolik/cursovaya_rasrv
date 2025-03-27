package usecase

import (
	"context"
	"httpserver/domain"
	"time"
)

type accountUsecase struct {
	userRepository    domain.UserRepository
	sessionRepository domain.SessionRepository
	contextTimeout    time.Duration
}

func NewAccountUsecase(userRepository domain.UserRepository, sessionRepository domain.SessionRepository, timeout time.Duration) domain.AccountUsecase {
	return &accountUsecase{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		contextTimeout:    timeout,
	}
}

func (au *accountUsecase) AccountBySession(ctx context.Context, session string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, au.contextTimeout)
	defer cancel()
	id, err := au.sessionRepository.Get(ctx, session)
	if err != nil {
		return nil, err
	}
	user, err := au.userRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (au *accountUsecase) Create(ctx context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, au.contextTimeout)
	defer cancel()
	return au.userRepository.Create(ctx, user)
}
