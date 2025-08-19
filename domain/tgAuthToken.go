package domain

import (
	"context"
	"errors"
	"time"
)

var ErrTokenGeneration = errors.New("failed to generate token")

type TGAuthToken struct {
	Token     string
	ExpiresAt time.Time
	ChatID    int64
}

type TGAuthTokenRepository interface {
	Create(ctx context.Context, telegramID int64) (*TGAuthToken, error)
	Get(ctx context.Context, token string) (*TGAuthToken, error)
	Apply(ctx context.Context, authToken *TGAuthToken, user *User) error
}
