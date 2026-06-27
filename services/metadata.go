package services

import (
	"errors"
	"regexp"

	"github.com/dotping-me/baastard/payload"
	"github.com/dotping-me/baastard/schema"
	"gorm.io/gorm"
)

type MetadataService struct {
	db *gorm.DB
}

// For validation purposes
var entityNameRegex = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)
var reservedEntities = map[string]bool{
	"users":    true,
	"sessions": true,
	"metadata": true,
	"fields":   true,
	"auth":     true,
	"system":   true,
	"admin":    true,
}

func NewMetadataService(db *gorm.DB) *MetadataService {
	return &MetadataService{db: db}
}

// Helper method that validates an Entity name
func isValidEntityName(name string) bool {
	return entityNameRegex.MatchString(name)
}

// Helper method that validates and allows only supported types
func isValidFieldType(t string) bool {
	switch t {
	case schema.FieldText:
		return true

	case schema.FieldBool:
		return true

	case schema.FieldDate:
		return true

	case schema.FieldJSON:
		return true

	case schema.FieldUUID:
		return true

	case schema.FieldEmail:
		return true

	case schema.FieldFloat:
		return true

	case schema.FieldInteger:
		return true

	case schema.FieldRelation:
		return true

	case schema.FieldFile:
		return true

	default:
		return false
	}
}

// Method that creates an EntityMetadata and respective FieldMetadata, with rollback safety
func (s *MetadataService) Create(req payload.CreateEntityMetadataRequest) (*schema.EntityMetadata, error) {

	var existing schema.EntityMetadata
	err := s.db.Where("name = ?", req.Name).First(&existing).Error
	if err == nil {
		return nil, errors.New("Entity already exists!")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Validating parsed name
	if !isValidEntityName(req.Name) {
		return nil, errors.New("Invalid Entity name!")
	}

	if reservedEntities[req.Name] {
		return nil, errors.New("Cannot use reserved Entity name!")
	}

	// Creates the Entity and stores it in the database
	em := schema.EntityMetadata{
		Name: req.Name,
	}

	// Because every field must be validated before storing everything in the
	// database. Starts a transaction which can safely rollback in the events of
	// failure.

	tx := s.db.Begin()

	err = tx.Create(&em).Error // Creates EntityMetadata and stores metadata
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Validating every field in that Entity
	duplicateFields := make(map[string]bool)

	for _, fm := range req.FieldsMetadata {
		if duplicateFields[fm.Name] {
			tx.Rollback()
			return nil, errors.New("Duplicate field name: " + fm.Name)
		}

		duplicateFields[fm.Name] = true

		if !isValidFieldType(fm.Type) {
			tx.Rollback()
			return nil, errors.New("Unsupported type: " + fm.Type)
		}

		// Creates the field and stores it
		fieldMetadata := schema.FieldMetadata{
			EntityMetadataID: em.ID, // EntityMetadata - FieldMetada relationship

			Name:     fm.Name,
			Type:     fm.Type,
			Required: fm.Required,
			Nullable: fm.Nullable,
		}

		err := tx.Create(&fieldMetadata).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Should be good. Now attempts to commit the transaction!
	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	// Builds response
	err = s.db.Preload("FieldsMetadata").First(&em, "id = ?", em.ID).Error
	if err != nil {
		return nil, err
	}

	return &em, nil
}

// Method that retrieves all Entities Metadata
func (s *MetadataService) GetAll() ([]schema.EntityMetadata, error) {
	var entitiesMetadata []schema.EntityMetadata

	err := s.db.Preload("FieldsMetadata").Find(&entitiesMetadata).Error
	if err != nil {
		return nil, err
	}

	return entitiesMetadata, nil
}

// Method that searches for an Entity Metadata by ID
func (s *MetadataService) GetById(id string) (*schema.EntityMetadata, error) {
	var em schema.EntityMetadata

	err := s.db.Preload("FieldsMetadata").Where("id = ?", id).First(&em).Error
	if err != nil {
		return nil, err
	}

	return &em, nil
}

// Method that returns an EntityMetadata bearing a given name
func (s *MetadataService) GetByName(name string) (*schema.EntityMetadata, error) {
	var em schema.EntityMetadata

	err := s.db.Preload("FieldsMetadata").Where("name = ?", name).First(&em).Error
	if err != nil {
		return nil, err
	}

	return &em, nil
}

// Method that deletes a collection with its associated fields using database constraints
func (s *MetadataService) Delete(id string) error {
	em, err := s.GetById(id)
	if err != nil {
		return err
	}

	if em.System {
		return errors.New("Cannot delete System-level Collections!")
	}

	return s.db.Delete(em).Error
}
