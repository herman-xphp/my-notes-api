package impl

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/herman-xphp/my-notes-api/internal/domain"
	"github.com/herman-xphp/my-notes-api/internal/dto"
	"github.com/herman-xphp/my-notes-api/internal/repository"
	"github.com/herman-xphp/my-notes-api/internal/service"
	"gorm.io/gorm"
)

type noteService struct {
	noteRepo repository.NoteRepository
	userRepo repository.UserRepository
}

// Ensure noteService implements service.NoteService
var _ service.NoteService = (*noteService)(nil)

// NewNoteService creates a new note service instance
func NewNoteService(
	noteRepo repository.NoteRepository,
	userRepo repository.UserRepository,
) service.NoteService {
	return &noteService{
		noteRepo: noteRepo,
		userRepo: userRepo,
	}
}

func (s *noteService) Create(ctx context.Context, userID uint, req dto.CreateNoteRequest) (*dto.NoteResponse, error) {
	// Verify user exists
	if _, err := s.userRepo.FindByID(ctx, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Create note
	note := &domain.Note{
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
	}

	if err := s.noteRepo.Create(ctx, note); err != nil {
		return nil, fmt.Errorf("failed to create note: %w", err)
	}

	// Build response
	return &dto.NoteResponse{
		ID:        note.ID,
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: note.CreatedAt.Format(time.RFC3339),
		UpdatedAt: note.UpdatedAt.Format(time.RFC3339),
	}, nil
}
func (s *noteService) GetByID(ctx context.Context, noteID, userID uint) (*dto.NoteResponse, error) {
	// Find note with ownership check
	note, err := s.noteRepo.FindByID(ctx, noteID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.ErrNoteNotFound
		}
		return nil, fmt.Errorf("failed to find note: %w", err)
	}

	// Build response
	return &dto.NoteResponse{
		ID:        note.ID,
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: note.CreatedAt.Format(time.RFC3339),
		UpdatedAt: note.UpdatedAt.Format(time.RFC3339),
	}, nil
}
func (s *noteService) GetAll(ctx context.Context, userID uint, req dto.NoteQueryRequest) (*dto.NoteListResponse, error) {
	// Set default values
	req.SetDefaults()

	// Build query params
	params := repository.NoteQueryParams{
		UserID:   userID,
		Page:     req.Page,
		PageSize: req.PageSize,
		Search:   req.Search,
		SortBy:   req.SortBy,
		SortDir:  req.SortDir,
	}

	// Get notes
	notes, total, err := s.noteRepo.FindAll(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to find notes: %w", err)
	}

	// Convert to response DTOs
	noteResponses := make([]dto.NoteResponse, len(notes))
	for i, note := range notes {
		noteResponses[i] = dto.NoteResponse{
			ID:        note.ID,
			Title:     note.Title,
			Content:   note.Content,
			CreatedAt: note.CreatedAt.Format(time.RFC3339),
			UpdatedAt: note.UpdatedAt.Format(time.RFC3339),
		}
	}

	// Calculate pagination
	totalPages := dto.CalculateTotalPages(total, req.PageSize)

	// Build response
	return &dto.NoteListResponse{
		Notes: noteResponses,
		Pagination: dto.PaginationMeta{
			Page:       req.Page,
			PageSize:   req.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}
func (s *noteService) Update(ctx context.Context, noteID, userID uint, req dto.UpdateNoteRequest) (*dto.NoteResponse, error) {
	// Find note with ownership check
	note, err := s.noteRepo.FindByID(ctx, noteID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.ErrNoteNotFound
		}
		return nil, fmt.Errorf("failed to find note: %w", err)
	}

	// Update fields if provided
	if req.Title != "" {
		note.Title = req.Title
	}
	if req.Content != "" {
		note.Content = req.Content
	}

	// Save changes
	if err := s.noteRepo.Update(ctx, note); err != nil {
		return nil, fmt.Errorf("failed to update note: %w", err)
	}

	// Build response
	return &dto.NoteResponse{
		ID:        note.ID,
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: note.CreatedAt.Format(time.RFC3339),
		UpdatedAt: note.UpdatedAt.Format(time.RFC3339),
	}, nil
}
func (s *noteService) Delete(ctx context.Context, noteID, userID uint) error {
	// Find note with ownership check
	_, err := s.noteRepo.FindByID(ctx, noteID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service.ErrNoteNotFound
		}
		return fmt.Errorf("failed to find note: %w", err)
	}

	// Soft delete
	if err := s.noteRepo.Delete(ctx, noteID, userID); err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	return nil
}
func (s *noteService) Restore(ctx context.Context, noteID, userID uint) (*dto.NoteResponse, error) {
	// Restore the note
	if err := s.noteRepo.Restore(ctx, noteID, userID); err != nil {
		return nil, fmt.Errorf("failed to restore note: %w", err)
	}

	// Get the restored note
	note, err := s.noteRepo.FindByID(ctx, noteID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find restored note: %w", err)
	}

	// Build response
	return &dto.NoteResponse{
		ID:        note.ID,
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: note.CreatedAt.Format(time.RFC3339),
		UpdatedAt: note.UpdatedAt.Format(time.RFC3339),
	}, nil
}
func (s *noteService) HardDelete(ctx context.Context, noteID, userID uint) error {
	// Verify note exists and belongs to user
	_, err := s.noteRepo.FindByID(ctx, noteID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service.ErrNoteNotFound
		}
		return fmt.Errorf("failed to find note: %w", err)
	}

	// Permanent delete
	if err := s.noteRepo.HardDelete(ctx, noteID, userID); err != nil {
		return fmt.Errorf("failed to permanently delete note: %w", err)
	}

	return nil
}
