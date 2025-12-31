package repository

import (
	"context"
	"database/sql"

	"github.com/TimX-21/auth-service-go/internal/apperror"
	"github.com/TimX-21/auth-service-go/internal/auth/model"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) getExecutor(c context.Context) DbInterface {
	if c != nil {
		if tx, ok := GetTx(c); ok {
			return tx
		}
	}
	return r.db
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, user model.User) (*model.User, error) {
	conn := r.getExecutor(ctx)

	query := "SELECT id, email, password, is_verified, created_at, updated_at, deleted_at FROM users WHERE email = $1"

	err := conn.QueryRowContext(ctx, query, user.Email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.IsVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.ErrUserNotFound
		}
		return nil, apperror.ErrDatabase
	}

	return &user, nil
}

func (r *AuthRepository) CreateUser(ctx context.Context, user model.User) error {
	conn := r.getExecutor(ctx)

	query := "INSERT INTO users (username, email, password, is_verified, created_at, updated_at) VALUES ($1, $2, $3, $4, NOW(), NOW())"

	_, err := conn.ExecContext(ctx, query,
		user.Username,
		user.Email,
		user.Password,
		user.IsVerified,
	)
	if err != nil {
		return apperror.ErrDatabase
	}

	return nil
}

func (r *AuthRepository) CreateOTP(ctx context.Context, otp model.PasswordResetOTP) error {
	conn := r.getExecutor(ctx)

	query := `
	INSERT INTO password_reset_otps
		(user_id, otp_hash, expires_at)
	VALUES ($1, $2, $3)
	`

	_, err := conn.ExecContext(
		ctx, query,
		otp.UserId,
		otp.OTP,
		otp.ExpiresAt,
	)
	if err != nil {
		return apperror.ErrDatabase
	}

	return nil
}

func (r *AuthRepository) GetLatestValidByUserId(ctx context.Context, user model.User) (*model.PasswordResetOTP, error) {
	conn := r.getExecutor(ctx)

	query := `
	SELECT id, user_id, otp_hash, expires_at, used_at, attempt_count, created_at
	FROM password_reset_otps
	WHERE user_id = $1
		AND used_at IS NULL
		AND expires_at > NOW()
	ORDER BY created_at DESC
	LIMIT 1
	`

	var result model.PasswordResetOTP
	err := conn.QueryRowContext(ctx, query, user.ID).Scan(
		&result.Id,
		&result.UserId,
		&result.OTP,
		&result.ExpiresAt,
		&result.UsedAt,
		&result.AttemptCount,
		&result.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, apperror.ErrDatabase
	}

	return &result, nil
}

func (r *AuthRepository) IncrementAttempt(ctx context.Context, otp model.PasswordResetOTP) error {
	conn := r.getExecutor(ctx)

	query := `
	UPDATE password_reset_otps
	SET attempt_count = attempt_count + 1
	WHERE id = $1
	`

	_, err := conn.ExecContext(ctx, query, otp.Id)
	if err != nil {
		return apperror.ErrDatabase
	}

	return nil
}

func (r *AuthRepository) MarkUsed(ctx context.Context, otp model.PasswordResetOTP) error {
	conn := r.getExecutor(ctx)

	query := `
	UPDATE password_reset_otps
	SET used_at = NOW()
	WHERE id = $1
	`

	_, err := conn.ExecContext(ctx, query, otp.Id)
	if err != nil {
		return apperror.ErrDatabase
	}

	return nil
}

func (r *AuthRepository) CountRecentByUserId(ctx context.Context, otp model.PasswordResetOTP) (int, error) {
	conn := r.getExecutor(ctx)

	query := `
	SELECT COUNT (*)
	FROM password_reset_otps
	WHERE user_id = $1
		AND created_at >= $2
	`

	var count int
	err := conn.QueryRowContext(ctx, query, otp.UserId, otp.CreatedAt).Scan(&count)
	if err != nil {
		return 0, apperror.ErrDatabase
	}

	return count, nil
}

func (r *AuthRepository) GetLatestByUserId(ctx context.Context, otp model.PasswordResetOTP) (*model.PasswordResetOTP, error) {
	conn := r.getExecutor(ctx)

	query := `
	SELECT id, user_id, otp_hash, expires_at, used_at, attempt_count, created_at
	FROM password_reset_otps
	WHERE user_id = $1
	ORDER BY created_at DESC
	LIMIT 1
	`

	var result model.PasswordResetOTP
	err := conn.QueryRowContext(ctx, query, otp.UserId).Scan(
		&result.Id,
		&result.UserId,
		&result.OTP,
		&result.ExpiresAt,
		&result.UsedAt,
		&result.AttemptCount,
		&result.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, apperror.ErrDatabase
	}

	return &result, nil
}

func (r *AuthRepository) CreateResetToken(ctx context.Context, token model.PasswordResetToken) error {
	conn := r.getExecutor(ctx)

	query := `
	INSERT INTO password_reset_tokens
		(user_id, token_hash, expires_at)
	VALUES ($1, $2, $3)
	`

	_, err := conn.ExecContext(
		ctx, query,
		token.UserId,
		token.Token,
		token.ExpiresAt,
	)
	if err != nil {
		return apperror.ErrDatabase
	}

	return nil
}

func (r *AuthRepository) GetValidByHash(ctx context.Context, token model.PasswordResetToken) (*model.PasswordResetToken, error) {
	conn := r.getExecutor(ctx)

	query := `
	SELECT id, user_id, token_hash, expires_at, used_at, created_at
	FROM password_reset_tokens
	WHERE token_hash = $1
		AND used_at IS NULL
		AND expires_at > NOW()
	LIMIT 1
	`

	var result model.PasswordResetToken
	err := conn.QueryRowContext(ctx, query, token.Token).Scan(
		&result.Id,
		&result.UserId,
		&result.Token,
		&result.ExpiresAt,
		&result.UsedAt,
		&result.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, apperror.ErrDatabase
	}

	return &result, nil
}

func (r *AuthRepository) MarkTokenUsed(ctx context.Context, token model.PasswordResetToken) error {
	conn := r.getExecutor(ctx)

	query := `
	UPDATE password_reset_tokens
	SET used_at = NOW()
	WHERE id = $1
	`

	_, err := conn.ExecContext(ctx, query, token.Id)
	if err != nil {
		return apperror.ErrDatabase
	}

	return nil

}

func (r *AuthRepository) UpdatePasswordHash(ctx context.Context, user model.User) error {
	conn := r.getExecutor(ctx)

	query := `
	UPDATE users
	SET password = $1, updated_at = NOW()
	WHERE id = $2
	`

	_, err := conn.ExecContext(ctx, query, user.Password, user.ID)
	if err != nil {
		return apperror.ErrDatabase
	}

	return nil
}

func (r *AuthRepository) RevokeActiveVerificationTokens(ctx context.Context, userId int64) error {
	conn := r.getExecutor(ctx)

	query := `
	UPDATE email_verification_tokens
	SET revoked_at = NOW()
	WHERE user_id = $1 AND used_at IS NULL AND revoked_at IS NULL
	`

	if _, err := conn.ExecContext(ctx, query, userId); err != nil {
		return apperror.ErrDatabase
	}
	return nil
}

func (r *AuthRepository) GetVerificationTokenByHash(ctx context.Context, tokenHash string) (*model.EmailVerificationToken, error) {
	conn := r.getExecutor(ctx)

	query := `
	SELECT user_id, expires_at, used_at, revoked_at
	FROM email_verification_tokens
	WHERE token_hash = $1
	`

	var result model.EmailVerificationToken
	err := conn.QueryRowContext(ctx, query, tokenHash).Scan(
		&result.UserId,
		&result.ExpiresAt,
		&result.UsedAt,
		&result.RevokedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, apperror.ErrDatabase
	}

	return &result, nil
}

func (r *AuthRepository) MarkVerificationTokenUsed(ctx context.Context, tokenHash string) error {
	conn := r.getExecutor(ctx)

	query := `
	UPDATE email_verification_tokens
	SET used_at = NOW()
	WHERE token_hash = $1 AND used_at IS NULL AND revoked_at IS NULL
	`

	if _, err := conn.ExecContext(ctx, query, tokenHash); err != nil {
		return apperror.ErrDatabase
	}

	return nil
}

func (r *AuthRepository) VerifyUserEmail(ctx context.Context, userId int64) error {
	conn := r.getExecutor(ctx)

	query := `
	UPDATE users
	SET is_verified = TRUE,
		email_verified_at = NOW()
	WHERE id = $1 AND is_verified = FALSE
	`

	if _, err := conn.ExecContext(ctx, query, userId); err != nil {
		return apperror.ErrDatabase
	}
	return nil
}
