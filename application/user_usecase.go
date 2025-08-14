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
	if err != nil {
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

func (uc *UserUseCase) LinkTelegram(ctx context.Context, userID int64, chatID int64) error {
	user, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}
	user.TelegramChatID = chatID
	return uc.repo.Update(ctx, user)
}

func (uc *UserUseCase) UpdateUserRole(ctx context.Context, operator *domain.User, targetUserID int64, newRole string) error {
	if !operator.Permissions.Has(domain.PermUserManage) {
		return errors.New("no permission to manage users")
	}

	targetUser, err := uc.repo.GetByID(ctx, targetUserID)
	if err != nil {
		return err
	}
	if targetUser == nil {
		return ErrUserNotFound
	}

	var newPerm domain.Permission
	switch newRole {
	case "guest":
		newPerm = domain.RoleGuest
	case "regular":
		newPerm = domain.RoleRegular
	case "privileged":
		newPerm = domain.RolePrivileged
	case "editor":
		newPerm = domain.RoleEditor
	case "admin":
		newPerm = domain.RoleAdmin
	default:
		return errors.New("invalid role specified")
	}

	targetUser.Permissions = newPerm
	return uc.repo.Update(ctx, targetUser)
}
