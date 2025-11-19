package mysql

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/herman-xphp/my-notes-api/internal/domain"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Migrate tables
	err = db.AutoMigrate(&domain.User{}, &domain.Note{}, &domain.RefreshToken{})
	require.NoError(t, err)

	return db
}

func TestUserRepositoryCreate(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &domain.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashedpassword",
	}

	err := repo.Create(ctx, user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
}

func TestUserRepositoryFindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create user
	user := &domain.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashedpassword",
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Find user
	found, err := repo.FindByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, found.Email)
	assert.Equal(t, user.Name, found.Name)
}

func TestUserRepositoryFindByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create user
	user := &domain.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashedpassword",
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Find by email
	found, err := repo.FindByEmail(ctx, "john@example.com")
	assert.NoError(t, err)
	assert.Equal(t, user.ID, found.ID)
	assert.Equal(t, user.Name, found.Name)
}

func TestUserRepositoryFindByEmailNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	_, err := repo.FindByEmail(ctx, "notfound@example.com")
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestUserRepositoryUpdate(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create user
	user := &domain.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashedpassword",
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Update user
	user.Name = "John Updated"
	err = repo.Update(ctx, user)
	assert.NoError(t, err)

	// Verify update
	found, err := repo.FindByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "John Updated", found.Name)
}

func TestUserRepositoryDelete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create user
	user := &domain.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashedpassword",
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Delete user
	err = repo.Delete(ctx, user.ID)
	assert.NoError(t, err)

	// Verify deletion (soft delete)
	_, err = repo.FindByID(ctx, user.ID)
	assert.Error(t, err)
}

func TestUserRepositoryExistsByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create user
	user := &domain.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashedpassword",
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Check exists
	exists, err := repo.ExistsByEmail(ctx, "john@example.com")
	assert.NoError(t, err)
	assert.True(t, exists)

	// Check non-existent
	exists, err = repo.ExistsByEmail(ctx, "notfound@example.com")
	assert.NoError(t, err)
	assert.False(t, exists)
}
