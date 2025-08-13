package repository

import (
	"context"
	"testing"

	"1litw/domain"

	"github.com/stretchr/testify/require"
)

// Helper function to create a user for tests that need a user ID.
func createTestUser(t *testing.T, repo domain.UserRepository, username string) *domain.User {
	ctx := context.Background()
	user := &domain.User{
		Username:     username,
		PasswordHash: "password",
		Permissions:  domain.RoleAdmin, // Use a default permission
	}
	userID, err := repo.Create(ctx, user)
	require.NoError(t, err)
	createdUser, err := repo.GetByUsername(ctx, username)
	require.NoError(t, err)
	require.NotNil(t, createdUser)
	createdUser.ID = userID
	return createdUser
}

func TestShortURLRepository(t *testing.T) {
	userRepo := NewUserRepository(testDB)
	urlRepo := NewShortURLRepository(testDB)
	ctx := context.Background()

	// Create a user to associate with the URLs
	testUser := createTestUser(t, userRepo, "urltester_repo")

	// 1. Test Create
	shortURL := &domain.ShortURL{
		UserID:      testUser.ID,
		OriginalURL: "https://example.com/long-url",
		ShortPath:   "randompath_repo",
	}
	_, err := urlRepo.Create(ctx, shortURL)
	require.NoError(t, err)

	// 2. Test GetByPath
	foundURL, err := urlRepo.GetByPath(ctx, "randompath_repo")
	require.NoError(t, err)
	require.NotNil(t, foundURL)
	require.Equal(t, shortURL.OriginalURL, foundURL.OriginalURL)
	require.Equal(t, testUser.ID, foundURL.UserID)
	require.NotZero(t, foundURL.ID)

	// 3. Test ListByUserID
	userURLs, err := urlRepo.ListByUserID(ctx, testUser.ID)
	require.NoError(t, err)
	require.Len(t, userURLs, 1)

	// 4. Test Delete
	err = urlRepo.Delete(ctx, foundURL.ID)
	require.NoError(t, err)

	deletedURL, err := urlRepo.GetByPath(ctx, "randompath_repo")
	require.Error(t, err) // Expect an error because it's deleted
	require.Nil(t, deletedURL)

	// Verify the user has no URLs left
	remainingURLs, err := urlRepo.ListByUserID(ctx, testUser.ID)
	require.NoError(t, err)
	require.Len(t, remainingURLs, 0)
}
