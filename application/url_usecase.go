package application

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"1litw/domain"
	"1litw/utils"
)

var (
	ErrInvalidURL           = errors.New("invalid original URL format or protocol")
	ErrPathReserved         = errors.New("the requested custom path is reserved")
	ErrPathTaken            = errors.New("the requested custom path is already taken")
	ErrNoPermission         = errors.New("user does not have permission for this action")
	ErrShortURLNotFound     = errors.New("short URL not found")
	ErrDeleteNotAllowed     = errors.New("user is not allowed to delete this short URL")
	ErrCustomPathNotAllowed = errors.New("user is not allowed to create a custom path with this format")
)

// ReservedPathsPattern defines a regex for paths that cannot be used for custom short URLs.
var ReservedPathsPattern = regexp.MustCompile(`^/(api|auth|admin|assets|static)/.*|/favicon.ico|/robots.txt$`)

type URLUseCase struct {
	urlRepo  domain.ShortURLRepository
	userRepo domain.UserRepository
	// clickSink will be used for async click recording
	// clickSink ClickSink
}

func NewURLUseCase(urlRepo domain.ShortURLRepository, userRepo domain.UserRepository) *URLUseCase {
	return &URLUseCase{
		urlRepo:  urlRepo,
		userRepo: userRepo,
	}
}

func (uc *URLUseCase) CreateShortURL(ctx context.Context, user *domain.User, originalURL, customPath string) (*domain.ShortURL, error) {
	// 1. Validate Original URL
	if !isValidURL(originalURL) {
		return nil, ErrInvalidURL
	}

	// 2. Determine User (handle anonymous)
	var userID int64
	if user == nil {
		anonUser, err := uc.userRepo.GetByUsername(ctx, "anonymous")
		if err != nil || anonUser == nil {
			return nil, fmt.Errorf("failed to find anonymous user: %w", err)
		}
		userID = anonUser.ID
	} else {
		userID = user.ID
	}

	// 3. Handle Path
	shortPath := customPath
	if shortPath == "" {
		// Generate a random path
		shortPath = utils.GenerateRandomString(6) // Assuming a util function
	} else {
		// Validate custom path
		if err := uc.validateCustomPath(ctx, user, shortPath); err != nil {
			return nil, err
		}
	}

	// 4. Check for uniqueness
	existing, err := uc.urlRepo.GetByPath(ctx, shortPath)
	if err != nil && err != domain.ErrNotFound {
		return nil, fmt.Errorf("failed to check path existence: %w", err)
	}
	if existing != nil {
		return nil, ErrPathTaken
	}

	// 5. Create and save the ShortURL
	newURL := &domain.ShortURL{
		ShortPath:   shortPath,
		OriginalURL: originalURL,
		UserID:      userID,
		CreatedAt:   time.Now(),
	}

	id, err := uc.urlRepo.Create(ctx, newURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create short URL: %w", err)
	}
	newURL.ID = id

	return newURL, nil
}

func (uc *URLUseCase) validateCustomPath(ctx context.Context, user *domain.User, path string) error {
	if ReservedPathsPattern.MatchString("/" + path) {
		return ErrPathReserved
	}

	// Guests cannot create custom paths
	if user == nil {
		return ErrNoPermission
	}

	// Path with prefix: @username/some-path
	if strings.HasPrefix(path, "@") {
		parts := strings.SplitN(path, "/", 2)
		usernameFromPath := strings.TrimPrefix(parts[0], "@")

		// Check if user has permission to create prefixed URLs
		if !user.Permissions.Has(domain.PermCreatePrefix) {
			return ErrCustomPathNotAllowed
		}
		// Check if the username in the path matches the user's own username
		if usernameFromPath != user.Username {
			return ErrCustomPathNotAllowed
		}
		return nil
	}

	// Path without prefix (any path)
	if !user.Permissions.Has(domain.PermCreateAny) {
		return ErrCustomPathNotAllowed
	}

	return nil
}

func (uc *URLUseCase) DeleteShortURLByID(ctx context.Context, user *domain.User, shortURLID int64) error {
	if user == nil {
		return ErrNoPermission // Guests can't delete
	}

	shortURL, err := uc.urlRepo.GetByID(ctx, shortURLID)
	if err != nil {
		return fmt.Errorf("failed to get short URL: %w", err)
	}
	if shortURL == nil {
		return ErrShortURLNotFound
	}

	// Check permissions
	isOwner := shortURL.UserID == user.ID
	canDeleteOwn := user.Permissions.Has(domain.PermDeleteOwn)
	canDeleteAny := user.Permissions.Has(domain.PermDeleteAny)

	if canDeleteAny {
		return uc.urlRepo.Delete(ctx, shortURL.ID)
	}

	if isOwner && canDeleteOwn {
		return uc.urlRepo.Delete(ctx, shortURL.ID)
	}

	return ErrDeleteNotAllowed
}

func (uc *URLUseCase) ListByUser(ctx context.Context, user *domain.User) ([]domain.ShortURL, error) {
	if user == nil {
		return nil, ErrNoPermission
	}
	return uc.urlRepo.ListByUserID(ctx, user.ID)
}

func (uc *URLUseCase) GetByPath(ctx context.Context, path string) (*domain.ShortURL, error) {
	return uc.urlRepo.GetByPath(ctx, path)
}

// isValidURL checks if a string is a valid URL with http or https protocol.
func isValidURL(rawURL string) bool {
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}
