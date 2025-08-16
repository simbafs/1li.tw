package application

import (
	"context"
	"errors"
	"time"

	"1litw/domain"
	"1litw/utils"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user with this username already exists")
	ErrPermissionDenied   = errors.New("permission denied")
)

type UserUseCase struct {
	repo      domain.UserRepository
	jwtSecret string
}

func NewUserUseCase(repo domain.UserRepository, jwtSecret string) *UserUseCase {
	return &UserUseCase{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (uc *UserUseCase) Register(ctx context.Context, username, password string) (*domain.User, error) {
	existing, err := uc.repo.GetByUsername(ctx, username)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUserExists
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username:     username,
		PasswordHash: hashedPassword,
		Permissions:  domain.RoleRegular, // Default role for new users
	}

	id, err := uc.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = id
	return user, nil
}

func (uc *UserUseCase) Login(ctx context.Context, username, password string) (string, *domain.User, error) {
	user, err := uc.repo.GetByUsername(ctx, username)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, ErrUserNotFound
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", nil, ErrInvalidCredentials
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"unm": user.Username,
		"prm": user.Permissions,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		return "", nil, err
	}

	return tokenString, user, nil
}

func (uc *UserUseCase) GetMe(ctx context.Context, userID int64) (*domain.User, error) {
	return uc.repo.GetByID(ctx, userID)
}

func (uc *UserUseCase) GetAnonymousUser(ctx context.Context) (*domain.User, error) {
	user, err := uc.repo.GetByUsername(ctx, "anonymous")
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, errors.New("critical: anonymous user not found in database")
		}
		return nil, err
	}
	return user, nil
}

func (uc *UserUseCase) GetUserByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
	return uc.repo.GetByTelegramID(ctx, telegramID)
}

// func (uc *UserUseCase) LinkTelegram(ctx context.Context, userID int64, chatID int64) error {
// 	user, err := uc.repo.GetByID(ctx, userID)
// 	if err != nil {
// 		return err
// 	}
// 	if user == nil {
// 		return ErrUserNotFound
// 	}
// 	user.TelegramChatID = chatID
// 	return uc.repo.Update(ctx, user)
// }

func (uc *UserUseCase) List(ctx context.Context, operator *domain.User) ([]*domain.User, error) {
	if !operator.Permissions.Has(domain.PermUserManage) {
		return nil, ErrPermissionDenied
	}

	users, err := uc.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		user.PasswordHash = ""
	}

	return users, nil
}

func (uc *UserUseCase) UpdatePermissions(ctx context.Context, operator *domain.User, targetID int64, permissions domain.Permission) error {
	if targetID == domain.AnonymousID {
		return ErrPermissionDenied
	}

	isSelf := operator.ID == targetID
	canManageOther := operator.Permissions.Has(domain.PermUserManage)

	if !isSelf && !canManageOther {
		return ErrPermissionDenied
	}

	return uc.repo.UpdatePermissions(ctx, targetID, permissions)
}

func (uc *UserUseCase) Delete(ctx context.Context, operator *domain.User, targetID int64) error {
	if targetID == domain.AnonymousID {
		return ErrPermissionDenied
	}

	isSelf := operator.ID == targetID
	canManageOther := operator.Permissions.Has(domain.PermUserManage)

	if !isSelf && !canManageOther {
		return ErrPermissionDenied
	}

	return uc.repo.Delete(ctx, targetID)
}
