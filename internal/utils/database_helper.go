package utils

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// DBHelper provides reusable database utilities
type DBHelper struct{}

// NewDBHelper creates a new database helper instance
func NewDBHelper() *DBHelper {
	return &DBHelper{}
}

// WithTimeout creates context with timeout for database operations
func (h *DBHelper) WithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout == 0 {
		timeout = 5 * time.Second // Default timeout
	}
	return context.WithTimeout(parent, timeout)
}

// Transaction executes function within a database transaction
func (h *DBHelper) Transaction(db *gorm.DB, fn func(*gorm.DB) error) error {
	return db.Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

// Exists checks if a record exists
func (h *DBHelper) Exists(db *gorm.DB, model interface{}, conditions ...interface{}) (bool, error) {
	var count int64
	err := db.Model(model).Where(conditions[0], conditions[1:]...).Count(&count).Error
	return count > 0, err
}

// FindByID finds record by ID with type safety
func (h *DBHelper) FindByID(db *gorm.DB, model interface{}, id uint) error {
	return db.First(model, id).Error
}

// SoftDelete performs soft delete
func (h *DBHelper) SoftDelete(db *gorm.DB, model interface{}, id uint) error {
	return db.Delete(model, id).Error
}

// HardDelete performs permanent delete
func (h *DBHelper) HardDelete(db *gorm.DB, model interface{}, id uint) error {
	return db.Unscoped().Delete(model, id).Error
}

// Restore restores soft-deleted record
func (h *DBHelper) Restore(db *gorm.DB, model interface{}, id uint) error {
	return db.Model(model).Unscoped().Where("id = ?", id).Update("deleted_at", nil).Error
}

// CountRecords counts records with conditions
func (h *DBHelper) CountRecords(db *gorm.DB, model interface{}, conditions ...interface{}) (int64, error) {
	var count int64
	query := db.Model(model)

	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}

	err := query.Count(&count).Error
	return count, err
}

// Paginate applies pagination to query
func (h *DBHelper) Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 {
			pageSize = 10
		}
		if pageSize > 100 {
			pageSize = 100
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
