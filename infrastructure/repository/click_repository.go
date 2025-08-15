package repository

import (
	"context"
	"database/sql"
	"time"

	"1litw/domain"
	"1litw/sqlc"
)

type clickRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

// NewClickRepository creates a new instance of ClickRepository.
func NewClickRepository(db *sql.DB) *clickRepository {
	return &clickRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (r *clickRepository) Insert(ctx context.Context, c *domain.URLClick) (int64, error) {
	id, err := r.queries.CreateURLClick(ctx, sqlc.CreateURLClickParams{
		ShortURLID:   c.ShortURLID,
		CountryCode:  sql.NullString{String: c.CountryCode, Valid: c.CountryCode != ""},
		OSName:       sql.NullString{String: c.OSName, Valid: c.OSName != ""},
		BrowserName:  sql.NullString{String: c.BrowserName, Valid: c.BrowserName != ""},
		RawUserAgent: sql.NullString{String: c.RawUserAgent, Valid: c.RawUserAgent != ""},
		IPAddress:    sql.NullString{String: c.IPAddress, Valid: c.IPAddress != ""},
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *clickRepository) CountByShortURLID(ctx context.Context, shortURLID int64) (int64, error) {
	return r.queries.CountClicksByShortURLID(ctx, shortURLID)
}

func (r *clickRepository) AggregateByTime(ctx context.Context, shortURLID int64, from, to time.Time) ([]domain.TimeBucketCount, error) {
	rows, err := r.queries.GetClickStatsByTime(ctx, sqlc.GetClickStatsByTimeParams{
		ShortURLID: shortURLID,
		From:       from,
		To:         to,
	})
	if err != nil {
		return nil, err
	}

	buckets := make([]domain.TimeBucketCount, len(rows))
	for i, row := range rows {
		// The time comes from SQLite as a string, so we need to parse it.
		timeStr, ok := row.TimeBucket.(string)
		if !ok {
			// Handle or log the error appropriately
			continue
		}
		t, err := time.Parse("2006-01-02T15:04:05Z", timeStr)
		if err != nil {
			// Handle or log the error appropriately
			continue
		}
		buckets[i] = domain.TimeBucketCount{
			BucketStart: t,
			Count:       row.Count,
		}
	}
	return buckets, nil
}

func (r *clickRepository) AggregateByCountry(ctx context.Context, shortURLID int64, from, to time.Time) ([]domain.KeyCount, error) {
	rows, err := r.queries.GetClickStatsByCountry(ctx, sqlc.GetClickStatsByCountryParams{
		ShortURLID: shortURLID,
		From:       from,
		To:         to,
	})
	if err != nil {
		return nil, err
	}

	counts := make([]domain.KeyCount, len(rows))
	for i, row := range rows {
		counts[i] = domain.KeyCount{
			Key:   row.CountryCode.String,
			Count: row.Count,
		}
	}
	return counts, nil
}

func (r *clickRepository) AggregateByOS(ctx context.Context, shortURLID int64, from, to time.Time) ([]domain.KeyCount, error) {
	rows, err := r.queries.GetClickStatsByOS(ctx, sqlc.GetClickStatsByOSParams{
		ShortURLID: shortURLID,
		From:       from,
		To:         to,
	})
	if err != nil {
		return nil, err
	}

	counts := make([]domain.KeyCount, len(rows))
	for i, row := range rows {
		counts[i] = domain.KeyCount{
			Key:   row.OSName.String,
			Count: row.Count,
		}
	}
	return counts, nil
}

func (r *clickRepository) AggregateByBrowser(ctx context.Context, shortURLID int64, from, to time.Time) ([]domain.KeyCount, error) {
	rows, err := r.queries.GetClickStatsByBrowser(ctx, sqlc.GetClickStatsByBrowserParams{
		ShortURLID: shortURLID,
		From:       from,
		To:         to,
	})
	if err != nil {
		return nil, err
	}

	counts := make([]domain.KeyCount, len(rows))
	for i, row := range rows {
		counts[i] = domain.KeyCount{
			Key:   row.BrowserName.String,
			Count: row.Count,
		}
	}
	return counts, nil
}

func (r *clickRepository) GetUnprocessedClicks(ctx context.Context, limit int64) ([]sqlc.GetUnprocessedClicksRow, error) {
	return r.queries.GetUnprocessedClicks(ctx, limit)
}

func (r *clickRepository) UpdateClickGeoInfo(ctx context.Context, arg sqlc.UpdateClickGeoInfoParams) error {
	return r.queries.UpdateClickGeoInfo(ctx, arg)
}
