package domain

import (
	"context"
	"time"
)

// URLClick represents a single click event on a short URL.
type URLClick struct {
	ID           int64
	ShortURLID   int64
	ClickedAt    time.Time
	CountryCode  string
	OSName       string
	BrowserName  string
	RawUserAgent string
	IPAddress    string
	Country      string
	RegionName   string
	City         string
	Lat          float64
	Lon          float64
	ISP          string
	ASInfo       string
	IsProcessed  bool
}

// TimeBucketCount is used for aggregating click counts over time intervals.
type TimeBucketCount struct {
	BucketStart time.Time `json:"bucketStart"`
	Count       int64     `json:"count"`
}

// KeyCount is a generic structure for aggregating counts by a string key (like country, OS, or browser).
type KeyCount struct {
	Key   string `json:"key"`
	Count int64  `json:"count"`
}

// ClickRepository defines the interface for accessing click analytics data.
type ClickRepository interface {
	Create(ctx context.Context, c *URLClick) (int64, error)
	CountByShortURLID(ctx context.Context, shortURLID int64) (int64, error)
	AggregateByTime(ctx context.Context, shortURLID int64, from, to time.Time) ([]TimeBucketCount, error)
	AggregateByCountry(ctx context.Context, shortURLID int64, from, to time.Time) ([]KeyCount, error)
	AggregateByOS(ctx context.Context, shortURLID int64, from, to time.Time) ([]KeyCount, error)
	AggregateByBrowser(ctx context.Context, shortURLID int64, from, to time.Time) ([]KeyCount, error)
}
