package dto

import "time"

type GetUserDataRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type GetUserDataResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
