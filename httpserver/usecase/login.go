package usecase

import (
	"context"
	"crypto/rand"
	"httpserver/domain"
	"math/big"
	"time"
)

type loginUsecase struct {
	userRepository    domain.UserRepository
	sessionRepository domain.SessionRepository
	contextTimeout    time.Duration
}

func NewLoginUsecase(userRepository domain.UserRepository, sessionRepository domain.SessionRepository, timeout time.Duration) domain.LoginUsecase {
	return &loginUsecase{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		contextTimeout:    timeout,
	}
}

func (lu *loginUsecase) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, lu.contextTimeout)
	defer cancel()
	return lu.userRepository.GetByEmail(ctx, email)
}

func (lu *loginUsecase) SaveSession(ctx context.Context, id int) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, lu.contextTimeout)
	defer cancel()
	session, err := GenerateRandomString(32)
	if err != nil {
		return "", err
	}
	lu.sessionRepository.Save(ctx, session, id)
	return session, nil
}

const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"

func GenerateRandomString(n int) (string, error) {
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}
	return string(ret), nil
}
