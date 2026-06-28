package services

import (
	"errors"
	"fmt"

	"github.com/dotping-me/baastard/schema"
	"gorm.io/gorm"
)

type CrudService struct {
	db *gorm.DB
}

func NewCrudService(db *gorm.DB) *CrudService {
	return &CrudService{db: db}
}

// Helper method that looks up a table configuration
func (s *CrudService) GetTableConfig(name string) (*schema.TableConfig, error) {

	var tableConfig schema.TableConfig
	err := s.db.Where("name = ?", name).First(&tableConfig).Error
	if err != nil {
		return nil, err
	}

	// Table is not meant to be exposed
	if tableConfig.Hidden {
		return nil, errors.New("Table is not set to be exposed!")
	}

	return &tableConfig, nil
}

// Makes sure that table is not read-only
func (s *CrudService) isWritable(name string) error {

	table, err := s.GetTableConfig(name)
	if err != nil {
		return err
	}

	if table.ReadOnly {
		return errors.New("Table is read-only!")
	}

	return nil
}

// Returns all rows from a table
func (s *CrudService) GetRows(tableName string) ([]map[string]any, error) {

	_, err := s.GetTableConfig(tableName)
	if err != nil {
		return nil, err
	}

	var rows []map[string]any
	err = s.db.Table(tableName).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// Returns a single row (by id) from a table. Note that ID is of type `any`
// because database entity schemas are created through SQL.
func (s *CrudService) GetRow(tableName string, id any) (map[string]any, error) {

	_, err := s.GetTableConfig(tableName)
	if err != nil {
		return nil, err
	}

	row := make(map[string]any)
	err = s.db.Table(tableName).Where("id = ?", id).Take(&row).Error // TODO: Drop GORM for Dynamic Tables. Use raw SQL execution.
	if err != nil {
		return nil, err
	}

	return row, nil
}

// Inserting data into a table. Note that data is of type `any`
// because database entity schemas are created through SQL.
func (s *CrudService) Insert(tableName string, data map[string]any) error {

	err := s.isWritable(tableName) // Checks if not read-only
	if err != nil {
		return err
	}

	return s.db.Table(tableName).Create(data).Error
}

// Updating data. These methods are here to facilitate exposing dynamic
// API endpoints for CRUD.
func (s *CrudService) Update(tableName string, id any, data map[string]any) error {

	err := s.isWritable(tableName)
	if err != nil {
		return err
	}

	return s.db.Table(tableName).Where("id = ?", id).Updates(data).Error
}

// Deletes data by id from a table. Deleting a table entirely would require SQL.
func (s *CrudService) Delete(tableName string, id any) error {

	err := s.isWritable(tableName)
	if err != nil {
		return err
	}

	return s.db.Table(tableName).Where("id = ?", id).Delete(nil).Error
}

// Helper method to check if a table exists.
func (s *CrudService) Exists(tableName string) bool {
	var count int64

	s.db.Raw(fmt.Sprintf("SELECT COUNT(*) FROM pragma_table_info('%s')",
		tableName)).Scan(&count)

	return count > 0
}

// Returns the column metadata of a table, hence returning a representation
// of that table.
func (s *CrudService) Columns(table string) ([]schema.ColumnMetadata, error) {

	var columns []schema.ColumnMetadata
	query := fmt.Sprintf("PRAGMA table_info(%s)", table)
	err := s.db.Raw(query).Scan(&columns).Error
	if err != nil {
		return nil, err
	}

	return columns, nil
}
