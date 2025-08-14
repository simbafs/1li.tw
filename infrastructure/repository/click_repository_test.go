package repository

import (
	"context"
	"testing"
	"time"

	"1litw/domain"

	"github.com/stretchr/testify/require"
)

func TestClickRepository(t *testing.T) {
	userRepo := NewUserRepository(testDB)
	urlRepo := NewShortURLRepository(testDB)
	clickRepo := NewClickRepository(testDB)
	ctx := context.Background()

	// Setup: Create a user and a short URL
	testUser := createTestUser(t, userRepo, "clicktester_repo")
	shortURL := &domain.ShortURL{
		UserID:      testUser.ID,
		OriginalURL: "https://example.com/for-clicking",
		ShortPath:   "clickpath_repo",
	}
	shortURLID, err := urlRepo.Create(ctx, shortURL)
	require.NoError(t, err)

	// 1. Test Insert Click
	click := &domain.URLClick{
		ShortURLID:   shortURLID,
		RawUserAgent: "Go-Test",
		CountryCode:  "TW",
		OSName:       "Linux",
		BrowserName:  "Go",
	}
	_, err = clickRepo.Insert(ctx, click)
	require.NoError(t, err)

	// 2. Test CountByShortURLID
	count, err := clickRepo.CountByShortURLID(ctx, shortURLID)
	require.NoError(t, err)
	require.Equal(t, int64(1), count)

	// 3. Test Aggregation functions
	// Use a wider time range to ensure the click is included.
	from := time.Now().Add(-24 * time.Hour)
	to := time.Now().Add(24 * time.Hour)

	// Test AggregateByTime
	timeBuckets, err := clickRepo.AggregateByTime(ctx, shortURLID, from, to)
	require.NoError(t, err)
	require.NotEmpty(t, timeBuckets)
	require.Equal(t, int64(1), timeBuckets[0].Count)

	// Test AggregateByCountry
	countryCounts, err := clickRepo.AggregateByCountry(ctx, shortURLID, from, to)
	require.NoError(t, err)
	require.NotEmpty(t, countryCounts)
	require.Equal(t, "TW", countryCounts[0].Key)
	require.Equal(t, int64(1), countryCounts[0].Count)

	// Test AggregateByOS
	osCounts, err := clickRepo.AggregateByOS(ctx, shortURLID, from, to)
	require.NoError(t, err)
	require.NotEmpty(t, osCounts)
	require.Equal(t, "Linux", osCounts[0].Key)
	require.Equal(t, int64(1), osCounts[0].Count)

	// Test AggregateByBrowser
	browserCounts, err := clickRepo.AggregateByBrowser(ctx, shortURLID, from, to)
	require.NoError(t, err)
	require.NotEmpty(t, browserCounts)
	require.Equal(t, "Go", browserCounts[0].Key)
	require.Equal(t, int64(1), browserCounts[0].Count)
}
