package utils

import (
	"time"

	"github.com/herman-xphp/my-notes-api/internal/domain"
	"github.com/herman-xphp/my-notes-api/internal/dto"
)

// ResponseHelper provides reusable response transformation utilities
type ResponseHelper struct{}

// NewResponseHelper creates a new response helper instance
func NewResponseHelper() *ResponseHelper {
	return &ResponseHelper{}
}

// ToUserResponse converts domain User to DTO UserResponse
func (h *ResponseHelper) ToUserResponse(user *domain.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}
}

// ToNoteResponse converts domain Note to DTO NoteResponse
func (h *ResponseHelper) ToNoteResponse(note *domain.Note) dto.NoteResponse {
	return dto.NoteResponse{
		ID:        note.ID,
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: note.CreatedAt.Format(time.RFC3339),
		UpdatedAt: note.UpdatedAt.Format(time.RFC3339),
	}
}

// ToNoteResponses converts slice of Notes to slice of NoteResponse
func (h *ResponseHelper) ToNoteResponses(notes []domain.Note) []dto.NoteResponse {
	responses := make([]dto.NoteResponse, len(notes))
	for i, note := range notes {
		responses[i] = h.ToNoteResponse(&note)
	}
	return responses
}

// BuildAuthResponse builds complete auth response with tokens
func (h *ResponseHelper) BuildAuthResponse(
	user *domain.User,
	accessToken string,
	refreshToken string,
	expiresIn int64,
) *dto.AuthResponse {
	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		User:         h.ToUserResponse(user),
	}
}

// BuildPaginationMeta builds pagination metadata
func (h *ResponseHelper) BuildPaginationMeta(
	page int,
	pageSize int,
	total int64,
) dto.PaginationMeta {
	totalPages := dto.CalculateTotalPages(total, pageSize)

	return dto.PaginationMeta{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	}
}

// BuildNoteListResponse builds paginated note list response
func (h *ResponseHelper) BuildNoteListResponse(
	notes []domain.Note,
	page int,
	pageSize int,
	total int64,
) *dto.NoteListResponse {
	return &dto.NoteListResponse{
		Notes:      h.ToNoteResponses(notes),
		Pagination: h.BuildPaginationMeta(page, pageSize, total),
	}
}
