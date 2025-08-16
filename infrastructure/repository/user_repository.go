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
	return toDomainUser(user), nil
}

func (r *userRepository) List(ctx context.Context) ([]*domain.User, error) {
	users, err := r.queries.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	var result []*domain.User
	for _, user := range users {
		result = append(result, toDomainUser(user))
	}
	return result, nil
}

func (r *userRepository) UpdateTelegramID(ctx context.Context, id int64, telegramID int64) error {
	return r.queries.UpdateUserTelegramID(ctx, sqlc.UpdateUserTelegramIDParams{
		ID:             id,
		TelegramChatID: sql.NullInt64{Int64: telegramID, Valid: true},
	})
}

func (r *userRepository) UpdatePermissions(ctx context.Context, userID int64, permissions domain.Permission) error {
	return r.queries.UpdateUserPermissions(ctx, sqlc.UpdateUserPermissionsParams{
		ID:          userID,
		Permissions: int64(permissions),
	})
}

func (r *userRepository) Delete(ctx context.Context, userID int64) error {
	return r.queries.DeleteUser(ctx, userID)
}

func toDomainUser(user sqlc.User) *domain.User {
	return &domain.User{
		ID:             user.ID,
		Username:       user.Username,
		PasswordHash:   user.PasswordHash,
		Permissions:    domain.Permission(user.Permissions),
		TelegramChatID: user.TelegramChatID.Int64,
		CreatedAt:      user.CreatedAt,
	}
}
