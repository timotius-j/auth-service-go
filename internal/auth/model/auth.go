package model

import "time"

type User struct {
	ID        int64
	Email     string
	Password  string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
