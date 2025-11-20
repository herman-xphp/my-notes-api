package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/herman-xphp/my-notes-api/internal/domain"
	"github.com/herman-xphp/my-notes-api/internal/dto"
	"github.com/herman-xphp/my-notes-api/internal/repository"
	repomock "github.com/herman-xphp/my-notes-api/internal/repository/mock"
)

func setupNoteServiceTest() (NoteService, *repomock.NoteRepository, *repomock.UserRepository) {
	mockNoteRepo := new(repomock.NoteRepository)
	mockUserRepo := new(repomock.UserRepository)
	service := NewNoteService(mockNoteRepo, mockUserRepo)
	return service, mockNoteRepo, mockUserRepo
}

func TestNoteService_Create_Success(t *testing.T) {
	service, mockNoteRepo, mockUserRepo := setupNoteServiceTest()
	ctx := context.Background()

	userID := uint(1)
	req := dto.CreateNoteRequest{
		Title:   "Test Note",
		Content: "This is a test note",
	}

	existingUser := &domain.User{
		ID:    userID,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// Mock expectations
	mockUserRepo.On("FindByID", ctx, userID).Return(existingUser, nil)
	mockNoteRepo.On("Create", ctx, mock.AnythingOfType("*domain.Note")).Return(nil)

	// Execute
	result, err := service.Create(ctx, userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Title, result.Title)
	assert.Equal(t, req.Content, result.Content)
	assert.NotZero(t, result.ID)
	assert.NotEmpty(t, result.CreatedAt)
	assert.NotEmpty(t, result.UpdatedAt)

	mockUserRepo.AssertExpectations(t)
	mockNoteRepo.AssertExpectations(t)
}

func TestNoteService_Create_UserNotFound(t *testing.T) {
	service, mockNoteRepo, mockUserRepo := setupNoteServiceTest()
	ctx := context.Background()

	userID := uint(999)
	req := dto.CreateNoteRequest{
		Title:   "Test Note",
		Content: "This is a test note",
	}

	// Mock expectations
	mockUserRepo.On("FindByID", ctx, userID).Return(nil, gorm.ErrRecordNotFound)

	// Execute
	result, err := service.Create(ctx, userID, req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrNoteNotFound, err)
	assert.Nil(t, result)

	mockUserRepo.AssertExpectations(t)
	mockNoteRepo.AssertNotCalled(t, "Create")
}

func TestNoteService_GetByID_Success(t *testing.T) {
	service, mockNoteRepo, _ := setupNoteServiceTest()
	ctx := context.Background()

	userID := uint(1)
	noteID := uint(1)

	existingNote := &domain.Note{
		ID:        noteID,
		UserID:    userID,
		Title:     "Test Note",
		Content:   "This is a test note",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock expectations
	mockNoteRepo.On("FindByID", ctx, noteID, userID).Return(existingNote, nil)

	// Execute
	result, err := service.GetByID(ctx, noteID, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, existingNote.ID, result.ID)
	assert.Equal(t, existingNote.Title, result.Title)
	assert.Equal(t, existingNote.Content, result.Content)

	mockNoteRepo.AssertExpectations(t)
}

func TestNoteService_GetByID_NotFound(t *testing.T) {
	service, mockNoteRepo, _ := setupNoteServiceTest()
	ctx := context.Background()

	userID := uint(1)
	noteID := uint(999)

	// Mock expectations
	mockNoteRepo.On("FindByID", ctx, noteID, userID).Return(nil, gorm.ErrRecordNotFound)

	// Execute
	result, err := service.GetByID(ctx, noteID, userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrNoteNotFound, err)
	assert.Nil(t, result)

	mockNoteRepo.AssertExpectations(t)
}

func TestNoteService_GetAll_Success(t *testing.T) {
	service, mockNoteRepo, _ := setupNoteServiceTest()
	ctx := context.Background()

	userID := uint(1)
	req := dto.NoteQueryRequest{
		Page:     1,
		PageSize: 10,
		Search:   "",
		SortBy:   "created_at",
		SortDir:  "desc",
	}

	notes := []domain.Note{
		{
			ID:        1,
			UserID:    userID,
			Title:     "Note 1",
			Content:   "Content 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			UserID:    userID,
			Title:     "Note 2",
			Content:   "Content 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Mock expectations
	mockNoteRepo.On("FindAll", ctx, mock.AnythingOfType("repository.NoteQueryParams")).
		Return(notes, int64(2), nil)

	// Execute
	result, err := service.GetAll(ctx, userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Notes, 2)
	assert.Equal(t, int64(2), result.Pagination.TotalItems)
	assert.Equal(t, 1, result.Pagination.Page)
	assert.Equal(t, 10, result.Pagination.PageSize)
	assert.Equal(t, 1, result.Pagination.TotalPages)

	mockNoteRepo.AssertExpectations(t)
}

func TestNoteService_GetAll_WithSearch(t *testing.T) {
	service, mockNoteRepo, _ := setupNoteServiceTest()
	ctx := context.Background()

	userID := uint(1)
	req := dto.NoteQueryRequest{
		Page:     1,
		PageSize: 10,
		Search:   "golang",
		SortBy:   "created_at",
		SortDir:  "desc",
	}

	notes := []domain.Note{
		{
			ID:        1,
			UserID:    userID,
			Title:     "Golang Tutorial",
			Content:   "Learn Go programming",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Mock expectations
	mockNoteRepo.On("FindAll", ctx, mock.MatchedBy(func(params repository.NoteQueryParams) bool {
		return params.Search == "golang"
	})).Return(notes, int64(1), nil)

	// Execute
	result, err := service.GetAll(ctx, userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Notes, 1)
	assert.Equal(t, "Golang Tutorial", result.Notes[0].Title)

	mockNoteRepo.AssertExpectations(t)
}

func TestNoteService_GetAll_EmptyResult(t *testing.T) {
	service, mockNoteRepo, _ := setupNoteServiceTest()
	ctx := context.Background()

	userID := uint(1)
	req := dto.NoteQueryRequest{
		Page:     1,
		PageSize: 10,
	}

	// Mock expectations
	mockNoteRepo.On("FindAll", ctx, mock.AnythingOfType("repository.NoteQueryParams")).
		Return([]domain.Note{}, int64(0), nil)

	// Execute
	result, err := service.GetAll(ctx, userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.Notes)
	assert.Equal(t, int64(0), result.Pagination.TotalItems)

	mockNoteRepo.AssertExpectations(t)
}

func TestNoteService_Update_Success(t *testing.T) {
	service, mockNoteRepo, _ := setupNoteServiceTest()
	ctx := context.Background()

	userID := uint(1)
	noteID := uint(1)

	existingNote := &domain.Note{
		ID:        noteID,
		UserID:    userID,
		Title:     "Original Title",
		Content:   "Original Content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	req := dto.UpdateNoteRequest{
		Title:   "Updated Title",
		Content: "Updated Content",
	}

	// Mock expectations
	mockNoteRepo.On("FindByID", ctx, noteID, userID).Return(existingNote, nil)
	mockNoteRepo.On("Update", ctx, mock.AnythingOfType("*domain.Note")).Return(nil)

	// Execute
	result, err := service.Update(ctx, noteID, userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Title, result.Title)
	assert.Equal(t, req.Content, result.Content)

	mockNoteRepo.AssertExpectations(t)
}

func TestNoteService_Update_PartialUpdate(t *testing.T) {
	service, mockNoteRepo, _ := setupNoteServiceTest()
	ctx := context.Background()

	userID := uint(1)
	noteID := uint(1)

	existingNote := &domain.Note{
		ID:        noteID,
		UserID:    userID,
		Title:     "Original Title",
		Content:   "Original Content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	req := dto.UpdateNoteRequest{
		Title: "Updated Title",
		// Content not updated
	}

	// Mock expectations
	mockNoteRepo.On("FindByID", ctx, noteID, userID).Return(existingNote, nil)
	mockNoteRepo.On("Update", ctx, mock.AnythingOfType("*domain.Note")).Return(nil)

	// Execute
	result, err := service.Update(ctx, noteID, userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Title, result.Title)
	assert.Equal(t, "Original Content", result.Content) // Should keep original

	mockNoteRepo.AssertExpectations(t)
}

func TestNoteService_Update_NotFound(t *testing.T) {
	service, mockNoteRepo, _ := setupNoteServiceTest()
	ctx := context.Background()

	userID := uint(1)
	noteID := uint(999)

	req := dto.UpdateNoteRequest{
		Title: "Updated Title",
	}

	// Mock expectations
	mockNoteRepo.On("FindByID", ctx, noteID, userID).Return(nil, gorm.ErrRecordNotFound)

	// Execute
	result, err := service.Update(ctx, noteID, userID, req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrNoteNotFound, err)
	assert.Nil(t, result)

	mockNoteRepo.AssertExpectations(t)
	mockNoteRepo.AssertNotCalled(t, "Update")
}

func TestNoteService_Delete_Success(t *testing.T) {
	service, mockNoteRepo, _ := setupNoteServiceTest()
	ctx := context.Background()

	userID := uint(1)
	noteID := uint(1)

	existingNote := &domain.Note{
		ID:     noteID,
		UserID: userID,
		Title:  "Test Note",
	}

	// Mock expectations
	mockNoteRepo.On("FindByID", ctx, noteID, userID).Return(existingNote, nil)
	mockNoteRepo.On("Delete", ctx, noteID, userID).Return(nil)

	// Execute
	err := service.Delete(ctx, noteID, userID)

	// Assert
	assert.NoError(t, err)
	mockNoteRepo.AssertExpectations(t)
}

func TestNoteService_Delete_NotFound(t *testing.T) {
	service, mockNoteRepo, _ := setupNoteServiceTest()
	ctx := context.Background()

	userID := uint(1)
	noteID := uint(999)

	// Mock expectations
	mockNoteRepo.On("FindByID", ctx, noteID, userID).Return(nil, gorm.ErrRecordNotFound)

	// Execute
	err := service.Delete(ctx, noteID, userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrNoteNotFound, err)

	mockNoteRepo.AssertExpectations(t)
	mockNoteRepo.AssertNotCalled(t, "Delete")
}

func TestNoteService_Restore_Success(t *testing.T) {
	service, mockNoteRepo, _ := setupNoteServiceTest()
	ctx := context.Background()

	userID := uint(1)
	noteID := uint(1)

	restoredNote := &domain.Note{
		ID:        noteID,
		UserID:    userID,
		Title:     "Restored Note",
		Content:   "This was deleted",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock expectations
	mockNoteRepo.On("Restore", ctx, noteID, userID).Return(nil)
	mockNoteRepo.On("FindByID", ctx, noteID, userID).Return(restoredNote, nil)

	// Execute
	result, err := service.Restore(ctx, noteID, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, restoredNote.Title, result.Title)

	mockNoteRepo.AssertExpectations(t)
}

func TestNoteService_HardDelete_Success(t *testing.T) {
	service, mockNoteRepo, _ := setupNoteServiceTest()
	ctx := context.Background()

	userID := uint(1)
	noteID := uint(1)

	existingNote := &domain.Note{
		ID:     noteID,
		UserID: userID,
		Title:  "To be permanently deleted",
	}

	// Mock expectations
	mockNoteRepo.On("FindByID", ctx, noteID, userID).Return(existingNote, nil)
	mockNoteRepo.On("HardDelete", ctx, noteID, userID).Return(nil)

	// Execute
	err := service.HardDelete(ctx, noteID, userID)

	// Assert
	assert.NoError(t, err)
	mockNoteRepo.AssertExpectations(t)
}

func TestNoteService_HardDelete_NotFound(t *testing.T) {
	service, mockNoteRepo, _ := setupNoteServiceTest()
	ctx := context.Background()

	userID := uint(1)
	noteID := uint(999)

	// Mock expectations
	mockNoteRepo.On("FindByID", ctx, noteID, userID).Return(nil, gorm.ErrRecordNotFound)

	// Execute
	err := service.HardDelete(ctx, noteID, userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrNoteNotFound, err)

	mockNoteRepo.AssertExpectations(t)
	mockNoteRepo.AssertNotCalled(t, "HardDelete")
}
