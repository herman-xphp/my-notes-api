package domain

import (
	"time"

	"gorm.io/gorm"
)

type Note struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index:idx_notes_user_id" json:"user_id"`
	Title     string         `gorm:"type:varchar(255);not null" json:"title"`
	Content   string         `gorm:"type:text" json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // Soft delete

	// Relations
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

// TableName specifies the table name for Note model
func (Note) TableName() string {
	return "notes"
}

// BeforeCreate hook
func (n *Note) BeforeCreate(tx *gorm.DB) error {
	return nil
}

// BeforeUpdate hook
func (n *Note) BeforeUpdate(tx *gorm.DB) error {
	return nil
}
