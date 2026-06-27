package schema

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Similar to schema.Collection, this one stores metadaba about the fields for
// that collection/database entity.

const (
	FieldText     = "text"
	FieldBool     = "bool"
	FieldDate     = "date"
	FieldJSON     = "json"
	FieldUUID     = "uuid"
	FieldEmail    = "email"
	FieldFloat    = "float"
	FieldInteger  = "integer"
	FieldRelation = "relation"
	FieldFile     = "file"
)

type FieldMetadata struct {
	ID               string `json:"id" gorm:"primaryKey"`
	EntityMetadataID string `json:"entityId" gorm:"not null;index"`

	Name     string `json:"name" gorm:"not null"`
	Type     string `json:"type" gorm:"not null"`
	Required bool   `json:"required"`
	Nullable bool   `json:"nullable"`
	Order    uint   `json:"order"`

	// TODO: Add Unique, Indexed and DefaultValue fields here!

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (fm *FieldMetadata) BeforeCreate(tx *gorm.DB) error {
	if fm.ID == "" {
		fm.ID = uuid.NewString()
	}

	fm.CreatedAt = time.Now()
	fm.UpdatedAt = time.Now()

	return nil
}

func (fm *FieldMetadata) BeforeUpdate(tx *gorm.DB) error {
	fm.UpdatedAt = time.Now()
	return nil
}
