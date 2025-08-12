package domain

import (
	"context"
	"time"
)

type User struct {
	ID           int64
	Username     string
	PasswordHash string
	Permissions  int
	TelegramID   int64
	CreatedAt    time.Time
}
type UserRepository interface {
	Create(ctx context.Context, user User) (int64, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	UpdateTelegramID(ctx context.Context, tgid int64) error
	UpdatePermissions(ctx context.Context, id int64, permissions int) error
}
