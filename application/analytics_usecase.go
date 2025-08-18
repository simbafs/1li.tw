package application

import (
	"context"
	"errors"
	"time"

	"1litw/domain"
)

var ErrAnalyticsNotFound = errors.New("no analytics data found for the given short URL")

// URLStats is a composite struct holding all analytics for a URL.
type URLStats struct {
	URL       *domain.ShortURL         `json:"url"`
	OwnerName string                   `json:"owner_name"`
	Total     int64                    `json:"total"`
	ByTime    []domain.TimeBucketCount `json:"by_time"`
	ByCountry []domain.KeyCount        `json:"by_country"`
	ByOS      []domain.KeyCount        `json:"by_os"`
	ByBrowser []domain.KeyCount        `json:"by_browser"`
}

type AnalyticsUseCase struct {
	clickRepo domain.ClickRepository
	urlRepo   domain.ShortURLRepository
}

func NewAnalyticsUseCase(clickRepo domain.ClickRepository, urlRepo domain.ShortURLRepository) *AnalyticsUseCase {
	return &AnalyticsUseCase{
		clickRepo: clickRepo,
		urlRepo:   urlRepo,
	}
}

func (a *AnalyticsUseCase) GetOverviewByID(ctx context.Context, user *domain.User, shortURLID int64, from, to time.Time) (*URLStats, error) {
	shortURL, err := a.urlRepo.GetByID(ctx, shortURLID)
	if err != nil {
		return nil, err
	}
	if shortURL == nil {
		return nil, ErrShortURLNotFound
	}

	// Permission Check
	isOwner := user != nil && shortURL.UserID == user.ID
	canViewOwn := user != nil && user.Permissions.Has(domain.PermViewOwnStats)
	canViewAny := user != nil && user.Permissions.Has(domain.PermViewAnyStats)

	if !(canViewAny || (isOwner && canViewOwn)) {
		return nil, ErrNoPermission
	}

	// Set default time range if not provided (e.g., last 30 days)
	if from.IsZero() || to.IsZero() {
		to = time.Now()
		from = to.AddDate(0, -1, 0) // Last 30 days
	}

	total, err := a.clickRepo.CountByShortURLID(ctx, shortURL.ID)
	if err != nil {
		return nil, err
	}

	byTime, err := a.clickRepo.AggregateByTime(ctx, shortURL.ID, from, to)
	if err != nil {
		return nil, err
	}

	byCountry, err := a.clickRepo.AggregateByCountry(ctx, shortURL.ID, from, to)
	if err != nil {
		return nil, err
	}

	byOS, err := a.clickRepo.AggregateByOS(ctx, shortURL.ID, from, to)
	if err != nil {
		return nil, err
	}

	byBrowser, err := a.clickRepo.AggregateByBrowser(ctx, shortURL.ID, from, to)
	if err != nil {
		return nil, err
	}

	stats := &URLStats{
		URL:       shortURL,
		OwnerName: user.Username,
		Total:     total,
		ByTime:    byTime,
		ByCountry: byCountry,
		ByOS:      byOS,
		ByBrowser: byBrowser,
	}

	return stats, nil
}
