package impl

import (
	"context"

	"github.com/herman-xphp/my-notes-api/internal/domain"
	"github.com/herman-xphp/my-notes-api/internal/dto"
	"github.com/herman-xphp/my-notes-api/internal/repository"
	"github.com/herman-xphp/my-notes-api/internal/service"
	"github.com/herman-xphp/my-notes-api/internal/utils"
)

type noteService struct {
	noteRepo       repository.NoteRepository
	userRepo       repository.UserRepository
	responseHelper *utils.ResponseHelper
	errorHelper    *utils.DBHelper
}

var _ service.NoteService = (*noteService)(nil)

func NewNoteService(
	noteRepo repository.NoteRepository,
	userRepo repository.UserRepository,
) service.NoteService {
	return &noteService{
		noteRepo:       noteRepo,
		userRepo:       userRepo,
		responseHelper: utils.NewResponseHelper(),
		errorHelper:    utils.NewDBHelper(),
	}
}

func (s *noteService) Create(ctx context.Context, userID uint, req dto.CreateNoteRequest) (*dto.NoteResponse, error) {
	// Verify user exists
	if _, err := s.userRepo.FindByID(ctx, userID); err != nil {
		// ✅ Reusable error handling
		return nil, utils.HandleDBError(err, service.ErrUserNotFound)
	}

	// Create note
	note := &domain.Note{
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
	}

	// ✅ Clean error wrapping
	if err := s.noteRepo.Create(ctx, note); err != nil {
		return nil, utils.WrapError("failed to create note", err)
	}

	// ✅ Reusable response transformation
	response := s.responseHelper.ToNoteResponse(note)
	return &response, nil
}

func (s *noteService) GetByID(ctx context.Context, noteID, userID uint) (*dto.NoteResponse, error) {
	note, err := s.noteRepo.FindByID(ctx, noteID, userID)
	if err != nil {
		// ✅ One liner
		return nil, utils.HandleDBError(err, service.ErrNoteNotFound)
	}

	response := s.responseHelper.ToNoteResponse(note)
	return &response, nil
}

func (s *noteService) GetAll(ctx context.Context, userID uint, req dto.NoteQueryRequest) (*dto.NoteListResponse, error) {
	req.SetDefaults()

	params := repository.NoteQueryParams{
		UserID:   userID,
		Page:     req.Page,
		PageSize: req.PageSize,
		Search:   req.Search,
		SortBy:   req.SortBy,
		SortDir:  req.SortDir,
	}

	notes, total, err := s.noteRepo.FindAll(ctx, params)
	if err != nil {
		return nil, utils.WrapError("failed to find notes", err)
	}

	// ✅ Reusable response builder
	response := s.responseHelper.BuildNoteListResponse(notes, req.Page, req.PageSize, total)
	return response, nil
}

func (s *noteService) Update(ctx context.Context, noteID, userID uint, req dto.UpdateNoteRequest) (*dto.NoteResponse, error) {
	note, err := s.noteRepo.FindByID(ctx, noteID, userID)
	if err != nil {
		return nil, utils.HandleDBError(err, service.ErrNoteNotFound)
	}

	// Update fields if provided
	if req.Title != "" {
		note.Title = req.Title
	}
	if req.Content != "" {
		note.Content = req.Content
	}

	if err := s.noteRepo.Update(ctx, note); err != nil {
		return nil, utils.WrapError("failed to update note", err)
	}

	response := s.responseHelper.ToNoteResponse(note)
	return &response, nil
}

func (s *noteService) Delete(ctx context.Context, noteID, userID uint) error {
	// Verify ownership
	if _, err := s.noteRepo.FindByID(ctx, noteID, userID); err != nil {
		return utils.HandleDBError(err, service.ErrNoteNotFound)
	}

	// ✅ Clean error wrapping
	return utils.CheckAndWrap(
		s.noteRepo.Delete(ctx, noteID, userID),
		"failed to delete note",
	)
}

func (s *noteService) Restore(ctx context.Context, noteID, userID uint) (*dto.NoteResponse, error) {
	if err := s.noteRepo.Restore(ctx, noteID, userID); err != nil {
		return nil, utils.WrapError("failed to restore note", err)
	}

	note, err := s.noteRepo.FindByID(ctx, noteID, userID)
	if err != nil {
		return nil, utils.WrapError("failed to find restored note", err)
	}

	response := s.responseHelper.ToNoteResponse(note)
	return &response, nil
}

func (s *noteService) HardDelete(ctx context.Context, noteID, userID uint) error {
	if _, err := s.noteRepo.FindByID(ctx, noteID, userID); err != nil {
		return utils.HandleDBError(err, service.ErrNoteNotFound)
	}

	return utils.CheckAndWrap(
		s.noteRepo.HardDelete(ctx, noteID, userID),
		"failed to permanently delete note",
	)
}
