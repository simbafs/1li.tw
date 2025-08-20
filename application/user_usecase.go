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
	ErrInvalidToken       = errors.New("invalid token")
)

type UserUseCase struct {
	jwtSecret       string
	userRepo        domain.UserRepository
	tgAuthTokenRepo domain.TGAuthTokenRepository
}

func NewUserUseCase(jwtSecret string, userRepo domain.UserRepository, tgAuthTokenRepo domain.TGAuthTokenRepository) *UserUseCase {
	return &UserUseCase{
		jwtSecret:       jwtSecret,
		userRepo:        userRepo,
		tgAuthTokenRepo: tgAuthTokenRepo,
	}
}

func (uc *UserUseCase) Register(ctx context.Context, username, password string) (*domain.User, error) {
	existing, err := uc.userRepo.GetByUsername(ctx, username)
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

	id, err := uc.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = id
	return user, nil
}

func (uc *UserUseCase) Login(ctx context.Context, username, password string) (string, *domain.User, error) {
	user, err := uc.userRepo.GetByUsername(ctx, username)
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
	return uc.userRepo.GetByID(ctx, userID)
}

func (uc *UserUseCase) GetAnonymousUser(ctx context.Context) (*domain.User, error) {
	user, err := uc.userRepo.GetByUsername(ctx, "anonymous")
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, errors.New("critical: anonymous user not found in database")
		}
		return nil, err
	}
	return user, nil
}

func (uc *UserUseCase) GetUserByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
	return uc.userRepo.GetByTelegramID(ctx, telegramID)
}

func (uc *UserUseCase) List(ctx context.Context, operator *domain.User) ([]*domain.User, error) {
	if !operator.Permissions.Has(domain.PermUserManage) {
		return nil, ErrPermissionDenied
	}

	users, err := uc.userRepo.List(ctx)
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

	return uc.userRepo.UpdatePermissions(ctx, targetID, permissions)
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

	return uc.userRepo.Delete(ctx, targetID)
}

func (uc *UserUseCase) PrepareLinkTelegram(ctx context.Context, telegramID int64) (string, error) {
	token, err := uc.tgAuthTokenRepo.Create(ctx, telegramID)
	if err != nil {
		return "", err
	}

	return token.Token, nil
}

func (uc *UserUseCase) LinkTelegram(ctx context.Context, token string, user *domain.User) error {
	authToken, err := uc.tgAuthTokenRepo.Get(ctx, token)
	if err != nil {
		return err
	}

	if authToken.ExpiresAt.Before(time.Now()) {
		return ErrInvalidToken
	}

	if err := uc.tgAuthTokenRepo.Apply(ctx, authToken, user); err != nil {
		return err
	}

	return nil
}
