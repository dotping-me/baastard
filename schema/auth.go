package schema

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Auth users are registered users on the system and store metadata such as
// emails and passwords, and are the entities tied to authentication logic and sessions
// (i.e. JWT)

type Auth struct {
	ID           string    `json:"id" 		gorm:"primaryKey"`
	Email        string    `json:"email"	gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-" 		gorm:"not null"`
	Verified     bool      `json:"verified"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// Automatically populates the ID field when creating a new auth user
func (a *Auth) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.NewString()
	}

	// Default values
	a.Verified = false
	a.CreatedAt = time.Now()

	return nil
}
