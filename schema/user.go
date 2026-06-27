package schema

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Auth users are registered users on the system and store metadata such as
// emails and passwords, and are the entities tied to authentication logic and sessions
// (i.e. JWT)

type User struct {
	ID           string    `json:"id" 		gorm:"primaryKey"`
	Email        string    `json:"email"	gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-" 		gorm:"not null"`
	Verified     bool      `json:"verified"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`

	TokenState uint `gorm:"default:0"` // Used to invalidate access tokens after logout
}

// Automatically populates the ID field when creating a new auth user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}

	// Default values
	u.Verified = false
	u.CreatedAt = time.Now()

	return nil
}
