package mock

import (
	"context"
	"time"

	"github.com/herman-xphp/my-notes-api/internal/domain"
	"github.com/herman-xphp/my-notes-api/internal/repository"
	"github.com/stretchr/testify/mock"
)

// UserRepository is a mock implementation of repository.UserRepository
type UserRepository struct {
	mock.Mock
}

// Ensure UserRepository implements repository.UserRepository interface
var _ repository.UserRepository = (*UserRepository)(nil)

func (m *UserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	// Set ID and timestamps for testing
	if args.Error(0) == nil {
		if user.ID == 0 {
			user.ID = 1
		}
		if user.CreatedAt.IsZero() {
			user.CreatedAt = time.Now()
		}
		if user.UpdatedAt.IsZero() {
			user.UpdatedAt = time.Now()
		}
	}

	return args.Error(0)
}

func (m *UserRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *UserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	if args.Error(0) == nil {
		user.UpdatedAt = time.Now()
	}
	return args.Error(0)
}

func (m *UserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}
