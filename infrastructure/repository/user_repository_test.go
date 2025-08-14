package repository

import (
	"context"
	"testing"

	"1litw/domain"

	"github.com/stretchr/testify/require"
)

func TestGetByUsername_NotFound(t *testing.T) {
	repo := NewUserRepository(testDB)
	ctx := context.Background()

	_, err := repo.GetByUsername(ctx, "nonexistentuser")
	require.Error(t, err)
	require.Equal(t, domain.ErrNotFound, err)
}

