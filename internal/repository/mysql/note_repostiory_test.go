package mysql

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/herman-xphp/my-notes-api/internal/domain"
	"github.com/herman-xphp/my-notes-api/internal/repository"
)

func createTestUser(t *testing.T, db *gorm.DB) *domain.User {
	user := &domain.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	err := db.Create(user).Error
	require.NoError(t, err)
	return user
}

func TestNoteRepositoryCreate(t *testing.T) {
	db := setupTestDB(t)
	repo := NewNoteRepository(db)
	ctx := context.Background()

	user := createTestUser(t, db)

	note := &domain.Note{
		UserID:  user.ID,
		Title:   "Test Note",
		Content: "This is a test note",
	}

	err := repo.Create(ctx, note)
	assert.NoError(t, err)
	assert.NotZero(t, note.ID)
}

func TestNoteRepositoryFindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewNoteRepository(db)
	ctx := context.Background()

	user := createTestUser(t, db)

	// Create note
	note := &domain.Note{
		UserID:  user.ID,
		Title:   "Test Note",
		Content: "This is a test note",
	}
	err := repo.Create(ctx, note)
	require.NoError(t, err)

	// Find note
	found, err := repo.FindByID(ctx, note.ID, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, note.Title, found.Title)
	assert.Equal(t, note.Content, found.Content)
}

func TestNoteRepositoryFindByIDWrongUser(t *testing.T) {
	db := setupTestDB(t)
	repo := NewNoteRepository(db)
	ctx := context.Background()

	user := createTestUser(t, db)

	// Create note
	note := &domain.Note{
		UserID:  user.ID,
		Title:   "Test Note",
		Content: "This is a test note",
	}
	err := repo.Create(ctx, note)
	require.NoError(t, err)

	// Try to find with wrong user ID
	_, err = repo.FindByID(ctx, note.ID, 999)
	assert.Error(t, err)
}

func TestNoteRepositoryFindAll(t *testing.T) {
	db := setupTestDB(t)
	repo := NewNoteRepository(db)
	ctx := context.Background()

	user := createTestUser(t, db)

	// Create multiple notes
	notes := []domain.Note{
		{UserID: user.ID, Title: "Note 1", Content: "Content 1"},
		{UserID: user.ID, Title: "Note 2", Content: "Content 2"},
		{UserID: user.ID, Title: "Note 3", Content: "Content 3"},
	}

	for i := range notes {
		err := repo.Create(ctx, &notes[i])
		require.NoError(t, err)
	}

	// Find all notes
	params := repository.NoteQueryParams{
		UserID:   user.ID,
		Page:     1,
		PageSize: 10,
	}
	result, total, err := repo.FindAll(ctx, params)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, result, 3)
}

func TestNoteRepositoryFindAllWithSearch(t *testing.T) {
	db := setupTestDB(t)
	repo := NewNoteRepository(db)
	ctx := context.Background()

	user := createTestUser(t, db)

	// Create notes with different content
	notes := []domain.Note{
		{UserID: user.ID, Title: "Golang Tutorial", Content: "Learn Go"},
		{UserID: user.ID, Title: "Python Guide", Content: "Learn Python"},
		{UserID: user.ID, Title: "JavaScript Tips", Content: "Learn JS"},
	}

	for i := range notes {
		err := repo.Create(ctx, &notes[i])
		require.NoError(t, err)
	}

	// Search for "Golang"
	params := repository.NoteQueryParams{
		UserID:   user.ID,
		Page:     1,
		PageSize: 10,
		Search:   "Golang",
	}
	result, total, err := repo.FindAll(ctx, params)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, result, 1)
	assert.Equal(t, "Golang Tutorial", result[0].Title)
}

func TestNoteRepositoryUpdate(t *testing.T) {
	db := setupTestDB(t)
	repo := NewNoteRepository(db)
	ctx := context.Background()

	user := createTestUser(t, db)

	// Create note
	note := &domain.Note{
		UserID:  user.ID,
		Title:   "Original Title",
		Content: "Original Content",
	}
	err := repo.Create(ctx, note)
	require.NoError(t, err)

	// Update note
	note.Title = "Updated Title"
	note.Content = "Updated Content"
	err = repo.Update(ctx, note)
	assert.NoError(t, err)

	// Verify update
	found, err := repo.FindByID(ctx, note.ID, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", found.Title)
	assert.Equal(t, "Updated Content", found.Content)
}

func TestNoteRepositoryDelete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewNoteRepository(db)
	ctx := context.Background()

	user := createTestUser(t, db)

	// Create note
	note := &domain.Note{
		UserID:  user.ID,
		Title:   "Test Note",
		Content: "To be deleted",
	}
	err := repo.Create(ctx, note)
	require.NoError(t, err)

	// Delete note (soft delete)
	err = repo.Delete(ctx, note.ID, user.ID)
	assert.NoError(t, err)

	// Verify deletion
	_, err = repo.FindByID(ctx, note.ID, user.ID)
	assert.Error(t, err)
}

func TestNoteRepositoryRestore(t *testing.T) {
	db := setupTestDB(t)
	repo := NewNoteRepository(db)
	ctx := context.Background()

	user := createTestUser(t, db)

	// Create and delete note
	note := &domain.Note{
		UserID:  user.ID,
		Title:   "Test Note",
		Content: "To be restored",
	}
	err := repo.Create(ctx, note)
	require.NoError(t, err)

	err = repo.Delete(ctx, note.ID, user.ID)
	require.NoError(t, err)

	// Restore note
	err = repo.Restore(ctx, note.ID, user.ID)
	assert.NoError(t, err)

	// Verify restoration
	found, err := repo.FindByID(ctx, note.ID, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, note.Title, found.Title)
}
