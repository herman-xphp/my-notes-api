package service

import (
	"context"
	"errors"

	"github.com/herman-xphp/my-notes-api/internal/dto"
)

// Common service errors
var (
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrUserNotFound        = errors.New("user not found")
	ErrNoteNotFound        = errors.New("note not found")
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
	ErrUnauthorizedAccess  = errors.New("unauthorized access to resource")
)

// AuthService defines authentication business logic interface
type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.AuthResponse, error)
	Logout(ctx context.Context, userID uint) error
	LogoutAll(ctx context.Context, userID uint) error
}

// NoteService defines note business logic interface
type NoteService interface {
	Create(ctx context.Context, userID uint, req dto.CreateNoteRequest) (*dto.NoteResponse, error)
	GetByID(ctx context.Context, noteID, userID uint) (*dto.NoteResponse, error)
	GetAll(ctx context.Context, userID uint, req dto.NoteQueryRequest) (*dto.NoteListResponse, error)
	Update(ctx context.Context, noteID, userID uint, req dto.UpdateNoteRequest) (*dto.NoteResponse, error)
	Delete(ctx context.Context, noteID, userID uint) error
	Restore(ctx context.Context, noteID, userID uint) (*dto.NoteResponse, error)
	HardDelete(ctx context.Context, noteID, userID uint) error
}

// Services is a container for all service interface
type Services struct {
	Auth AuthService
	Note NoteService
}
