package mysql

import (
	"context"

	"github.com/herman-xphp/my-notes-api/internal/domain"
	"github.com/herman-xphp/my-notes-api/internal/repository"
	"gorm.io/gorm"
)

type noteRepository struct {
	db *gorm.DB
}

// NewNoteRepository creates a new MySQL note repository instance
func NewNoteRepository(db *gorm.DB) repository.NoteRepository {
	return &noteRepository{db: db}
}

func (r *noteRepository) Create(ctx context.Context, note *domain.Note) error {
	return r.db.WithContext(ctx).Create(note).Error
}

func (r *noteRepository) FindByID(ctx context.Context, id, userID uint) (*domain.Note, error) {
	var note domain.Note
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Take(&note).Error
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *noteRepository) FindAll(ctx context.Context, params repository.NoteQueryParams) ([]domain.Note, int64, error) {
	var notes []domain.Note
	var total int64

	// Build base query
	query := r.db.WithContext(ctx).Model(&domain.Note{}).Where("user_id = ?", params.UserID)

	// Apply search filter
	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("title LIKE ? OR content LIKE ?", searchPattern, searchPattern)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortDir := params.SortDir
	if sortDir == "" {
		sortDir = "desc"
	}
	orderClause := sortBy + " " + sortDir

	// Apply pagination
	offset := (params.Page - 1) * params.PageSize
	err := query.Order(orderClause).Limit(params.PageSize).Offset(offset).Find(&notes).Error

	return notes, total, err
}

func (r *noteRepository) Update(ctx context.Context, note *domain.Note) error {
	return r.db.WithContext(ctx).Save(note).Error
}

func (r *noteRepository) Delete(ctx context.Context, id, userID uint) error {
	// Soft delete
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&domain.Note{}).Error
}

func (r *noteRepository) HardDelete(ctx context.Context, id, userID uint) error {
	// Permanent delete
	return r.db.WithContext(ctx).Unscoped().Where("id = ? AND user_id = ?", id, userID).Delete(&domain.Note{}).Error
}

func (r *noteRepository) Restore(ctx context.Context, id, userID uint) error {
	return r.db.WithContext(ctx).Model(&domain.Note{}).Unscoped().Where("id = ? AND user_id = ?", id, userID).Update("deleted_at", nil).Error
}
