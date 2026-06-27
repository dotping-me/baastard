package schema

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	ID           string `gorm:"primaryKey"`
	UserID       string `gorm:"index;not null"`
	RefreshToken string `gorm:"uniqueIndex;not null"`

	// Device Binding to prevent a possibly stolen refresh token to work on
	// another machine
	DeviceID  string `gorm:"index"`
	UserAgent string
	IPAddress string

	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *Session) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.NewString()
	}

	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()

	return nil
}

func (s *Session) BeforeUpdate(tx *gorm.DB) error {
	s.UpdatedAt = time.Now()
	return nil
}
