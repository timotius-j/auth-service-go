package service

import (
	"context"
	"time"

	"github.com/TimX-21/auth-service-go/internal/apperror"
	"github.com/TimX-21/auth-service-go/internal/auth/model"
	"github.com/TimX-21/auth-service-go/internal/auth/repository"
	"github.com/TimX-21/auth-service-go/internal/config"
	"github.com/TimX-21/auth-service-go/internal/util"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepository repository.AuthRepositoryItf
	txManager      repository.TransactionManager
	cfg            config.ResetConfig
	emailSender    util.DummyEmailSender
}

func NewAuthService(
	authRepository repository.AuthRepositoryItf,
	txManager repository.TransactionManager,
	cfg config.ResetConfig,
	emailSender *util.DummyEmailSender,
) *AuthService {
	return &AuthService{
		authRepository: authRepository,
		txManager:      txManager,
		cfg:            cfg,
		emailSender:    *emailSender,
	}
}

func (s *AuthService) GetUserDataService(ctx context.Context, user model.User) (*model.User, error) {
	userData, err := s.authRepository.GetUserByEmail(ctx, user)
	if err != nil {
		return nil, err
	}
	return userData, nil
}

func (s *AuthService) LoginService(ctx context.Context, user model.User) (string, error) {

	DbUserData, err := s.authRepository.GetUserByEmail(ctx, user)
	if err != nil {
		return "", apperror.ErrUserNotFound
	}

	InputPassword := user.Password
	DbPassword := DbUserData.Password

	err = bcrypt.CompareHashAndPassword([]byte(DbPassword), []byte(InputPassword))
	if err != nil {
		return "", apperror.ErrUnauthorized
	}

	// jwtsecret, err := util.GetJWTSecret()
	// if err != nil {
	// 	return "", apperror.ErrInternalServer
	// }

	token, err := util.GenerateJWT(user, false)
	if err != nil {
		return "", apperror.ErrInternalServer
	}

	return token, nil
}

func (s *AuthService) RegisterService(ctx context.Context, user model.User) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return apperror.ErrInternalServer
	}

	user.Password = string(hashedPassword)

	err = s.txManager.Do(ctx, func(txContext context.Context) error {

		err = s.authRepository.CreateUser(txContext, user)
		if err != nil {
			return err
		}

		return nil

	})

	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ForgotPasswordRequestService(ctx context.Context, user model.User) error {

	DbUserData, err := s.authRepository.GetUserByEmail(ctx, user)
	if err != nil {
		return err
	}

	if DbUserData == nil {
		return apperror.ErrUserNotFound
	}

	otpLookUp := model.PasswordResetOTP{
		UserId: DbUserData.ID,
	}

	last, err := s.authRepository.GetLatestByUserId(ctx, otpLookUp)
	if err != nil {
		return err
	}

	if last != nil && time.Since(last.CreatedAt) < s.cfg.ResendCooldown {
		return nil
	}

	since := time.Now().Add(-1 * time.Hour)

	otpRecord := model.PasswordResetOTP{
		UserId:    DbUserData.ID,
		CreatedAt: since,
	}

	count, err := s.authRepository.CountRecentByUserId(ctx, otpRecord)
	if err != nil {
		return err
	}

	if count >= s.cfg.ResendLimitPerHour {
		return apperror.ErrTooManyRequests
	}

	otp, err := util.GenerateOTP()
	if err != nil {
		return apperror.ErrInternalServer
	}

	err = s.txManager.Do(ctx, func(txContext context.Context) error {

		otpHash := util.HashWithHMAC(otp, s.cfg.HMACSecret)
		expiresAt := time.Now().Add(s.cfg.OTPTTL)

		err = s.authRepository.CreateOTP(txContext, model.PasswordResetOTP{
			UserId:    DbUserData.ID,
			OTP:       otpHash,
			ExpiresAt: expiresAt,
		})

		if err != nil {
			return err
		}

		return nil

	})

	if err != nil {
		return err
	}

	_ = s.emailSender.SendResetOTP(ctx, DbUserData.Email, otp)

	return nil
}

func (s *AuthService) VerifyResetOTPService(ctx context.Context, userInput model.User, otpInput model.PasswordResetOTP) (string, error) {
	userRecord, err := s.authRepository.GetUserByEmail(ctx, userInput)
	if err != nil {
		return "", err
	}
	if userRecord == nil {
		return "", apperror.ErrInvalidOTP
	}

	otpRow, err := s.authRepository.GetLatestValidByUserId(ctx, *userRecord)
	if err != nil {
		return "", err
	}
	if otpRow == nil {
		return "", apperror.ErrInvalidOTP
	}

	var rawToken string

	err = s.txManager.Do(ctx, func(txContext context.Context) error {

		if otpRow.AttemptCount >= s.cfg.MaxAttempts {
			_ = s.authRepository.MarkUsed(txContext, model.PasswordResetOTP{Id: otpRow.Id})
			return apperror.ErrInvalidOTP
		}

		inputHash := util.HashWithHMAC(otpInput.OTP, s.cfg.HMACSecret)
		if inputHash != otpRow.OTP {
			_ = s.authRepository.IncrementAttempt(txContext, model.PasswordResetOTP{Id: otpRow.Id})

			if otpRow.AttemptCount+1 >= s.cfg.MaxAttempts {
				_ = s.authRepository.MarkUsed(txContext, model.PasswordResetOTP{Id: otpRow.Id})
			}

			return apperror.ErrInvalidOTP
		}

		if err := s.authRepository.MarkUsed(txContext, model.PasswordResetOTP{Id: otpRow.Id}); err != nil {
			return err
		}

		rawToken, err = util.GenerateResetToken()
		if err != nil {
			return err
		}

		tokenHash := util.HashWithHMAC(rawToken, s.cfg.HMACSecret)
		tokenExpires := time.Now().Add(s.cfg.ResetTokenTTL)

		if err := s.authRepository.CreateResetToken(txContext, model.PasswordResetToken{
			UserId:    userRecord.ID,
			Token:     tokenHash,
			ExpiresAt: tokenExpires,
		}); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return rawToken, nil
}

func (s *AuthService) ResetPasswordService(ctx context.Context, tokenInput model.PasswordResetToken, userInput model.User) error {
	tokenHash := util.HashWithHMAC(tokenInput.Token, s.cfg.HMACSecret)

	tokenRow, err := s.authRepository.GetValidByHash(ctx, model.PasswordResetToken{
		Token: tokenHash,
	})
	if err != nil {
		return err
	}
	if tokenRow == nil {
		return apperror.ErrInvalidResetToken
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	if err != nil {
		return apperror.ErrInternalServer
	}

	return s.txManager.Do(ctx, func(txContext context.Context) error {
		if err := s.authRepository.UpdatePasswordHash(txContext, model.User{
			ID:       tokenRow.UserId,
			Password: string(newHashedPassword),
		}); err != nil {
			return err
		}

		if err := s.authRepository.MarkTokenUsed(txContext, model.PasswordResetToken{Id: tokenRow.Id}); err != nil {
			return err
		}

		return nil
	})
}

// func (s *AuthService) SendVerification(ctx context.Context, )
