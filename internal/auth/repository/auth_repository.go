package repository

import (
	"context"
	"database/sql"

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
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
