package domain

import (
	"context"
	"time"
)

type URLClick struct {
	ID          int64
	ShortURLID  int64
	ClickedAt   time.Time
	CountryCode *string
	OSName      *string
	BrowserName *string
	DeviceType  *string
	UserAgent   *string
}

type TimeBucketCount struct {
	BucketStart time.Time
	Count       int64
}

type KeyCount struct {
	Key   string
	Count int64
}

type ClickRepository interface {
	Insert(ctx context.Context, c URLClick) (int64, error)
	CountByShortURL(ctx context.Context, shortURLID int64) (int64, error)
	AggregateByTimeRange(ctx context.Context, shortURLID int64, from, to time.Time, bucket string) ([]TimeBucketCount, error)
	AggregateByCountry(ctx context.Context, shortURLID int64, from, to time.Time) ([]KeyCount, error)
	AggregateByOS(ctx context.Context, shortURLID int64, from, to time.Time) ([]KeyCount, error)
	AggregateByBrowser(ctx context.Context, shortURLID int64, from, to time.Time) ([]KeyCount, error)
}
