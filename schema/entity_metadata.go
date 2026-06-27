package schema

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// This schema basically stores metadata about user created collections/database
// entities. This whole architecture is basically a table to store tables:
// One table to rule them all!

type EntityMetadata struct {
	ID     string `json:"id" gorm:"primaryKey"`
	Name   string `json:"name" gorm:"uniqueIndex;not null"`
	System bool   `json:"system"` // System-specific: Cannot be changed

	FieldsMetadata []FieldMetadata `json:"fields,omitempty" gorm:"foreignKey:EntityMetadataID;constraint:OnDelete:CASCADE;"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (em *EntityMetadata) BeforeCreate(tx *gorm.DB) error {
	if em.ID == "" {
		em.ID = uuid.NewString()
	}

	em.CreatedAt = time.Now()
	em.UpdatedAt = time.Now()

	return nil
}

func (em *EntityMetadata) BeforeUpdate(tx *gorm.DB) error {
	em.UpdatedAt = time.Now()
	return nil
}
