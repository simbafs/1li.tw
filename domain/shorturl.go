package domain

import (
	"context"
	"time"
)

type ShortURL struct {
	ID          int64
	ShortPath   string
	OriginalURL string
	UserID      int64
	CreatedAt   time.Time
}
type ShortURLRepository interface {
	Create(ctx context.Context, shortURL ShortURL) (int64, error)
	GetByPath(ctx context.Context, path string) (*ShortURL, error)
	Delete(ctx context.Context, id int64) error
	ListByUser(ctx context.Context, userID int64) ([]ShortURL, error)
}
