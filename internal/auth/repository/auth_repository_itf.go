package repository

import (
	"context"

	"github.com/TimX-21/auth-service-go/internal/auth/model"
)

type AuthRepositoryItf interface {
	GetUserByEmail(ctx context.Context, user model.User) (*model.User, error)
}
