package mock

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/herman-xphp/my-notes-api/internal/domain"
	"github.com/herman-xphp/my-notes-api/internal/repository"
)

// NoteRepository is a mock implementation of repository.NoteRepository
type NoteRepository struct {
	mock.Mock
}

// Ensure NoteRepository implements repository.NoteRepository interface
var _ repository.NoteRepository = (*NoteRepository)(nil)

func (m *NoteRepository) Create(ctx context.Context, note *domain.Note) error {
	args := m.Called(ctx, note)
	// Set ID and timestamps for testing
	if args.Error(0) == nil {
		if note.ID == 0 {
			note.ID = 1
		}
		if note.CreatedAt.IsZero() {
			note.CreatedAt = time.Now()
		}
		if note.UpdatedAt.IsZero() {
			note.UpdatedAt = time.Now()
		}
	}
	return args.Error(0)
}

func (m *NoteRepository) FindByID(ctx context.Context, id, userID uint) (*domain.Note, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Note), args.Error(1)
}

func (m *NoteRepository) FindAll(ctx context.Context, params repository.NoteQueryParams) ([]domain.Note, int64, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]domain.Note), args.Get(1).(int64), args.Error(2)
}

func (m *NoteRepository) Update(ctx context.Context, note *domain.Note) error {
	args := m.Called(ctx, note)
	if args.Error(0) == nil {
		note.UpdatedAt = time.Now()
	}
	return args.Error(0)
}

func (m *NoteRepository) Delete(ctx context.Context, id, userID uint) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *NoteRepository) HardDelete(ctx context.Context, id, userID uint) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *NoteRepository) Restore(ctx context.Context, id, userID uint) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}
