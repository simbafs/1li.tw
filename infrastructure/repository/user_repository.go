package repository

import (
	"context"
	"database/sql"
	"errors"

	"1litw/domain"
	"1litw/sqlc"
)

type userRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) (int64, error) {
	createdUser, err := r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Permissions:  int64(user.Permissions),
	})
	if err != nil {
		return 0, err
	}
	return createdUser.ID, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	user, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &domain.User{
		ID:             user.ID,
		Username:       user.Username,
		PasswordHash:   user.PasswordHash,
		Permissions:    domain.Permission(user.Permissions),
		TelegramChatID: user.TelegramChatID.Int64,
		CreatedAt:      user.CreatedAt,
	}, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	user, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &domain.User{
		ID:             user.ID,
		Username:       user.Username,
		PasswordHash:   user.PasswordHash,
		Permissions:    domain.Permission(user.Permissions),
		TelegramChatID: user.TelegramChatID.Int64,
		CreatedAt:      user.CreatedAt,
	}, nil
}

func (r *userRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
	user, err := r.queries.GetUserByTelegramID(ctx, sql.NullInt64{Int64: telegramID, Valid: true})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &domain.User{
		ID:             user.ID,
		Username:       user.Username,
		PasswordHash:   user.PasswordHash,
		Permissions:    domain.Permission(user.Permissions),
		TelegramChatID: user.TelegramChatID.Int64,
		CreatedAt:      user.CreatedAt,
	}, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	// The current domain interface is simple. We only support updating telegram ID and permissions.
	// A more complex implementation might need to update other fields.
	if user.TelegramChatID != 0 {
		err := r.queries.UpdateUserTelegramID(ctx, sqlc.UpdateUserTelegramIDParams{
			ID:             user.ID,
			TelegramChatID: sql.NullInt64{Int64: user.TelegramChatID, Valid: true},
		})
		if err != nil {
			return err
		}
	}

	err := r.queries.UpdateUserPermissions(ctx, sqlc.UpdateUserPermissionsParams{
		ID:          user.ID,
		Permissions: int64(user.Permissions),
	})
	return err
}
