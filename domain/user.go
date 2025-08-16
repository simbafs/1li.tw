package domain

import (
	"context"
	"errors"
	"time"
)

// ErrNotFound is a common error for when an entity is not found.
var ErrNotFound = errors.New("not found")

// User represents a user in the system.
type User struct {
	ID             int64
	Username       string
	PasswordHash   string
	Permissions    Permission
	TelegramChatID int64
	CreatedAt      time.Time
}

// UserRepository defines the interface for user data operations.
// This interface is implemented by the infrastructure layer (e.g., using sqlc).
type UserRepository interface {
	Create(ctx context.Context, user *User) (int64, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByTelegramID(ctx context.Context, telegramID int64) (*User, error)
	List(ctx context.Context) ([]*User, error)
	UpdateTelegramID(ctx context.Context, id int64, telegramID int64) error
	UpdatePermissions(ctx context.Context, userID int64, permissions Permission) error
	Delete(ctx context.Context, userID int64) error
}
