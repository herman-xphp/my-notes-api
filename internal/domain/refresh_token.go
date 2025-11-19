package domain

import "time"

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index:idx_refresh_tokens_user_id" json:"user_id"`
	Token     string    `gorm:"type:varchar(500);uniqueIndex;index:idx_refresh_tokens_token;not null" json:"token"`
	ExpiresAt time.Time `gorm:"not null;index:idx_refresh_tokens_expires_at" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

// TableName specifies the table name for RefreshToken model
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// IsExpired checks if the refresh token has expired
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

// IsValid checks if the refresh token is still valid
func (rt *RefreshToken) IsValid() bool {
	return !rt.IsExpired()
}
