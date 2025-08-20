package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"1litw/domain"
	"1litw/sqlc"
)

var _ domain.ShortURLRepository = (*shortURLRepository)(nil)

type shortURLRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

// NewShortURLRepository creates a new instance of ShortURLRepository.
func NewShortURLRepository(db *sql.DB) domain.ShortURLRepository {
	return &shortURLRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (r *shortURLRepository) Create(ctx context.Context, shortURL *domain.ShortURL) (int64, error) {
	created, err := r.queries.CreateShortURL(ctx, sqlc.CreateShortURLParams{
		ShortPath:   shortURL.ShortPath,
		OriginalURL: shortURL.OriginalURL,
		UserID:      shortURL.UserID,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create short URL: %w", err)
	}
	return created.ID, nil
}

func (r *shortURLRepository) GetByPath(ctx context.Context, path string) (*domain.ShortURL, error) {
	url, err := r.queries.GetShortURLByPath(ctx, path)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get short URL by path: %w", err)
	}
	return &domain.ShortURL{
		ID:          url.ID,
		OriginalURL: url.OriginalURL,
		UserID:      url.UserID,
		ShortPath:   path, // Add the path back to the struct
	}, nil
}

func (r *shortURLRepository) GetByID(ctx context.Context, id int64) (*domain.ShortURL, error) {
	url, err := r.queries.GetShortURLByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get short URL by ID: %w", err)
	}
	return &domain.ShortURL{
		ID:          url.ID,
		ShortPath:   url.ShortPath,
		OriginalURL: url.OriginalURL,
		UserID:      url.UserID,
	}, nil
}

func (r *shortURLRepository) Delete(ctx context.Context, id int64) error {
	return r.queries.DeleteShortURL(ctx, id)
}

func (r *shortURLRepository) ListByUserID(ctx context.Context, userID int64) ([]domain.ShortURL, error) {
	rows, err := r.queries.ListShortURLsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list short URLs by user ID: %w", err)
	}
	urls := make([]domain.ShortURL, len(rows))
	for i, row := range rows {
		urls[i] = domain.ShortURL{
			ID:          row.ShortUrl.ID,
			ShortPath:   row.ShortUrl.ShortPath,
			OriginalURL: row.ShortUrl.OriginalURL,
			CreatedAt:   row.ShortUrl.CreatedAt,
			UserID:      userID,
			TotalClicks: row.TotalClicks,
		}
	}
	return urls, nil
}

func (r *shortURLRepository) ListAll(ctx context.Context) ([]domain.ShortURL, error) {
	rows, err := r.queries.ListAllShortURLs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list all short URLs: %w", err)
	}
	// We need a new struct to hold the result from ListAllShortURLs, let's define it here
	// This is not ideal, but for now it's the quickest way.
	// A better approach would be to define this in the domain layer if it's a common concept.
	type ShortURLWithUsername struct {
		domain.ShortURL
		OwnerUsername string
	}

	urls := make([]domain.ShortURL, len(rows))
	for i, row := range rows {
		urls[i] = domain.ShortURL{
			ID:          row.ShortUrl.ID,
			ShortPath:   row.ShortUrl.ShortPath,
			OriginalURL: row.ShortUrl.OriginalURL,
			CreatedAt:   row.ShortUrl.CreatedAt,
			UserID:      row.ShortUrl.UserID,
			TotalClicks: row.TotalClicks,
		}
	}
	return urls, nil
}

func (r *shortURLRepository) ListAllURLsWithUser(ctx context.Context) ([]domain.ShortURLWithUser, error) {
	rows, err := r.queries.ListAllURLsWithUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list all URLs with user: %w", err)
	}

	urls := make([]domain.ShortURLWithUser, len(rows))
	for i, row := range rows {
		urls[i] = domain.ShortURLWithUser{
			ShortURL: domain.ShortURL{
				ID:          row.ID,
				ShortPath:   row.ShortPath,
				OriginalURL: row.OriginalURL,
				UserID:      row.UserID,
				CreatedAt:   row.CreatedAt,
			},
			Username: row.Username,
		}
	}
	return urls, nil
}
