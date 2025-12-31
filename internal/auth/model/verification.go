package model

import "time"

type EmailVerificationToken struct {
	ID        int64
	UserId    int64
	Token     string
	ExpiresAt time.Time
	UsedAt    *time.Time
	RevokedAt *time.Time
}
