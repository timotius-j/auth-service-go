package service

import (
	"context"

	"github.com/TimX-21/auth-service-go/internal/auth/model"
	"github.com/TimX-21/auth-service-go/internal/auth/repository"
)

type AuthService struct {
	authRepository repository.AuthRepositoryItf
	txManager      repository.TransactionManager
}

func NewAuthService(authRepository repository.AuthRepositoryItf, txManager repository.TransactionManager) *AuthService {
	return &AuthService{
		authRepository: authRepository,
		txManager:      txManager,
	}
}

func (s *AuthService) GetUserDataService(ctx context.Context, user model.User) (*model.User, error) {
	userData, err := s.authRepository.GetUserByEmail(ctx, user)
	if err != nil {
		return nil, err
	}
	return userData, nil
}
