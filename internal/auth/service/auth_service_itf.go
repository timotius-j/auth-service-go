package service

import (
	"context"

	"github.com/TimX-21/auth-service-go/internal/auth/model"
)

type AuthServiceItf interface {
	GetUserDataService(ctx context.Context, user model.User) (*model.User, error)
}
