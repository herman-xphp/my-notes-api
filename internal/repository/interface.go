package repository

import (
	"context"

	"github.com/herman-xphp/my-notes-api/internal/domain"
)

// UserRepository defines method for user data access
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id uint) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uint) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

// NoteQueryParams holds parameters for querying notes
type NoteQueryParams struct {
	UserID   uint
	Page     int
	PageSize int
	Search   string
	SortBy   string
	SortDir  string
}

// NoteRepository defines methods for note data access
type NoteRepository interface {
	Create(ctx context.Context, note *domain.Note) error
	FindByID(ctx context.Context, id, userID uint) (*domain.Note, error)
	FindAll(ctx context.Context, params NoteQueryParams) ([]domain.Note, int64, error)
	Update(ctx context.Context, note *domain.Note) error
	Delete(ctx context.Context, id, userID uint) error
	HardDelete(ctx context.Context, id, userID uint) error
	Restore(ctx context.Context, id, userID uint) error
}

// RefreshTokenRepository defines methods for refresh token data access
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *domain.RefreshToken) error
	FindByToken(ctx context.Context, token string) (*domain.RefreshToken, error)
	FindByUserID(ctx context.Context, userID uint) ([]domain.RefreshToken, error)
	Delete(ctx context.Context, token string) error
	DeleteByUserID(ctx context.Context, userID uint) error
	DeleteExpired(ctx context.Context) error
}

// Repositories is a container for all repositories
// This makes it easier to inject all repos at once
type Repositories struct {
	User         UserRepository
	Note         NoteRepository
	RefreshToken RefreshTokenRepository
}
