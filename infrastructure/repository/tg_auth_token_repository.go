package repository

import (
	"context"
	"crypto/rand"
	"database/sql"
	"time"

	"1litw/domain"
	"1litw/sqlc"
)

var _ domain.TGAuthTokenRepository = (*tgAuthTokenRepository)(nil)

type tgAuthTokenRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

func NewTGAuthTokenRepository(db *sql.DB) *tgAuthTokenRepository {
	return &tgAuthTokenRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (t *tgAuthTokenRepository) Create(ctx context.Context, telegramID int64) (*domain.TGAuthToken, error) {
	token := rand.Text()
	if token == "" {
		return nil, domain.ErrTokenGeneration
	}

	expiresAt := time.Now().Add(10 * time.Minute)

	authToken, err := t.queries.CreateTelegramAuthToken(ctx, sqlc.CreateTelegramAuthTokenParams{
		Token:          token,
		ExpiresAt:      expiresAt,
		TelegramChatID: telegramID,
	})
	if err != nil {
		return nil, err
	}

	return toDomainTGAuthToken(authToken), nil
}

func (t *tgAuthTokenRepository) Get(ctx context.Context, token string) (*domain.TGAuthToken, error) {
	authToken, err := t.queries.GetTelegramAuthToken(ctx, token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return toDomainTGAuthToken(authToken), nil
}

func (t *tgAuthTokenRepository) Apply(ctx context.Context, authToken *domain.TGAuthToken, user *domain.User) error {
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := t.queries.WithTx(tx)

	err = qtx.UpdateUserTelegramID(ctx, sqlc.UpdateUserTelegramIDParams{
		ID: user.ID,
		TelegramChatID: sql.NullInt64{
			Int64: authToken.ChatID,
			Valid: true,
		},
	})
	if err != nil {
		return err
	}

	err = qtx.DeleteTelegramAuthToken(ctx, authToken.Token)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func toDomainTGAuthToken(authToken sqlc.TelegramAuthToken) *domain.TGAuthToken {
	return &domain.TGAuthToken{
		Token:     authToken.Token,
		ExpiresAt: authToken.ExpiresAt,
		ChatID:    authToken.TelegramChatID,
	}
}
