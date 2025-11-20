package mock

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/herman-xphp/my-notes-api/internal/domain"
	"github.com/herman-xphp/my-notes-api/internal/repository"
)

// RefreshTokenRepository is a mock implementation of repository.RefreshTokenRepository
type RefreshTokenRepository struct {
	mock.Mock
}

// Ensure RefreshTokenRepository implements repository.RefreshTokenRepository interface
var _ repository.RefreshTokenRepository = (*RefreshTokenRepository)(nil)

func (m *RefreshTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	args := m.Called(ctx, token)
	// Set ID and timestamp for testing
	if args.Error(0) == nil {
		if token.ID == 0 {
			token.ID = 1
		}
		if token.CreatedAt.IsZero() {
			token.CreatedAt = time.Now()
		}
	}
	return args.Error(0)
}

func (m *RefreshTokenRepository) FindByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RefreshToken), args.Error(1)
}

func (m *RefreshTokenRepository) FindByUserID(ctx context.Context, userID uint) ([]domain.RefreshToken, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.RefreshToken), args.Error(1)
}

func (m *RefreshTokenRepository) Delete(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *RefreshTokenRepository) DeleteByUserID(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *RefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
