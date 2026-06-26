package schema

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	ID           string `gorm:"primaryKey"`
	AuthID       string `gorm:"index;not null"`
	RefreshToken string `gorm:"uniqueIndex;not null"`

	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *Session) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.NewString()
	}

	return nil
}
