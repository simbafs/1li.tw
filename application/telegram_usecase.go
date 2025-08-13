package application

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"1litw/domain"
	"1litw/sqlc"
)

var (
	ErrTokenNotFound = errors.New("auth token not found")
	ErrTokenExpired  = errors.New("auth token has expired")
)

type TelegramUseCase struct {
	queries  *sqlc.Queries
	userRepo domain.UserRepository
}

func NewTelegramUseCase(db *sql.DB, userRepo domain.UserRepository) *TelegramUseCase {
	return &TelegramUseCase{
		queries:  sqlc.New(db),
		userRepo: userRepo,
	}
}

func (uc *TelegramUseCase) VerifyAndLink(ctx context.Context, token string, userID int64) error {
	// 1. Get the token from the database
	authToken, err := uc.queries.GetTelegramAuthToken(ctx, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrTokenNotFound
		}
		return err
	}

	// 2. Delete the token immediately to prevent reuse
	defer uc.queries.DeleteTelegramAuthToken(ctx, token)

	// 3. Check if the token is expired
	if time.Now().After(authToken.ExpiresAt) {
		return ErrTokenExpired
	}

	// 4. Get the user and update their Telegram ID
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	user.TelegramChatID = authToken.TelegramChatID
	return uc.userRepo.Update(ctx, user)
}
