package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type URL struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	ShortCode   string         `gorm:"uniqueIndex;not null;size:10" json:"short_code"`
	OriginalURL string         `gorm:"not null" json:"original_url"`
	UserID      *uuid.UUID     `gorm:"type:uuid;index" json:"user_id,omitempty"`
	ClickCount  int64          `gorm:"default:0" json:"click_count"`
	ExpiresAt   *time.Time     `json:"expires_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *URL) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

type CreateURLRequest struct {
	OriginalURL string `json:"original_url" binding:"required,url"`
	CustomCode  string `json:"custom_code,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"` // hours
}

type URLResponse struct {
	ID          uuid.UUID  `json:"id"`
	ShortCode   string     `json:"short_code"`
	ShortURL    string     `json:"short_url"`
	OriginalURL string     `json:"original_url"`
	ClickCount  int64      `json:"click_count"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

type ClickEvent struct {
	ShortCode string    `json:"short_code"`
	UserAgent string    `json:"user_agent"`
	IP        string    `json:"ip"`
	Referer   string    `json:"referer"`
	Timestamp time.Time `json:"timestamp"`
}
