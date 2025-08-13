package domain

import (
	"context"
	"time"
)

// ShortURL represents the core entity for a shortened URL.
type ShortURL struct {
	ID          int64
	ShortPath   string
	OriginalURL string
	UserID      int64
	CreatedAt   time.Time
	TotalClicks int64 // Added for presentation/API purposes
}

// ShortURLRepository defines the interface for short URL data operations.
type ShortURLRepository interface {
	Create(ctx context.Context, shortURL *ShortURL) (int64, error)
	GetByPath(ctx context.Context, path string) (*ShortURL, error)
	GetByID(ctx context.Context, id int64) (*ShortURL, error)
	Delete(ctx context.Context, id int64) error
	ListByUserID(ctx context.Context, userID int64) ([]ShortURL, error)
	ListAll(ctx context.Context) ([]ShortURL, error)
}
