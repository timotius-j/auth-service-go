package repository

import (
	"context"

	"github.com/TimX-21/auth-service-go/internal/auth/model"
)

type AuthRepositoryItf interface {
	GetUserByEmail(ctx context.Context, user model.User) (*model.User, error)
	CreateUser(ctx context.Context, user model.User) error
	CreateOTP(ctx context.Context, otp model.PasswordResetOTP) error
	GetLatestValidByUserId(ctx context.Context, user model.User) (*model.PasswordResetOTP, error)
	IncrementAttempt(ctx context.Context, otp model.PasswordResetOTP) error
	MarkUsed(ctx context.Context, otp model.PasswordResetOTP) error
	CountRecentByUserId(ctx context.Context, otp model.PasswordResetOTP) (int, error)
	GetLatestByUserId(ctx context.Context, otp model.PasswordResetOTP) (*model.PasswordResetOTP, error)
	CreateResetToken(ctx context.Context, token model.PasswordResetToken) error
	GetValidByHash(ctx context.Context, token model.PasswordResetToken) (*model.PasswordResetToken, error)
	MarkTokenUsed(ctx context.Context, token model.PasswordResetToken) error
	UpdatePasswordHash(ctx context.Context, user model.User) error
	RevokeActiveVerificationTokens(ctx context.Context, userId int64) error
	GetVerificationTokenByHash(ctx context.Context, tokenHash string) (*model.EmailVerificationToken, error)
	MarkVerificationTokenUsed(ctx context.Context, tokenHash string) error
	VerifyUserEmail(ctx context.Context, userId int64) error
}
