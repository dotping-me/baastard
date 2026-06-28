package services

import (
	"database/sql"
	"errors"

	"github.com/dotping-me/baastard/schema"
	"github.com/dotping-me/baastard/utils"
)

type CrudService struct {
	db           *sql.DB
	tableConfigs *TableConfigService
}

func NewCrudService(db *sql.DB, cfg *TableConfigService) *CrudService {
	return &CrudService{db: db, tableConfigs: cfg}
}

// Helper method that looks up a table configuration
func (s *CrudService) GetTableConfig(name string) (*schema.TableConfig, error) {
	return s.tableConfigs.Get(name)
}

// Makes sure that table is not read-only
func (s *CrudService) isWritable(name string) error {
	table, err := s.GetTableConfig(name)
	if err != nil {
		return err
	}

	if table.ReadOnly {
		return errors.New("Table is read-only!") // TODO: Investigate this because this error is returned twice!
	}

	return nil
}

// Returns all rows from a table
func (s *CrudService) GetRows(table string) ([]map[string]any, error) {
	if _, err := s.GetTableConfig(table); err != nil {
		return nil, err
	}

	query := utils.Select(table)
	rows, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	return utils.ScanRows(rows)
}

// Returns a single row (by id) from a table. Note that ID is of type `any`
// because database entity schemas are created through SQL.
func (s *CrudService) GetRow(table string, id any) (map[string]any, error) {
	if _, err := s.GetTableConfig(table); err != nil {
		return nil, err
	}

	query := utils.SelectByID(table)
	rows, err := s.db.Query(query, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results, err := utils.ScanRows(rows)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, sql.ErrNoRows
	}

	return results[0], nil
}

// Inserting data into a table. Note that data is of type `any`
// because database entity schemas are created through SQL.
func (s *CrudService) Insert(table string, data map[string]any) error {
	if err := s.isWritable(table); err != nil {
		return err
	}

	query, values := utils.Insert(table, data)
	_, err := s.db.Exec(query, values...)

	return err
}

// Updating data. These methods are here to facilitate exposing dynamic
// API endpoints for CRUD.
func (s *CrudService) Update(table string, id any, data map[string]any) error {
	if err := s.isWritable(table); err != nil {
		return err
	}

	query, args := utils.Update(table, id, data)
	_, err := s.db.Exec(query, args...)

	return err
}

// Deletes data by id from a table. Deleting a table entirely would require SQL.
func (s *CrudService) Delete(table string, id any) error {
	if err := s.isWritable(table); err != nil {
		return err
	}

	query := utils.Delete(table)
	_, err := s.db.Exec(query, id)

	return err
}

// Helper method to check if a table exists.
func (s *CrudService) Exists(table string) bool {
	query := utils.Exists()

	var count int
	err := s.db.QueryRow(query, table).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

// Returns the column metadata of a table, hence returning a representation
// of that table.
func (s *CrudService) Columns(table string) ([]schema.ColumnMetadata, error) {
	query := utils.Columns(table)

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var columns []schema.ColumnMetadata
	for rows.Next() {
		var c schema.ColumnMetadata

		err := rows.Scan(&c.CID, &c.Name, &c.Type, &c.NotNull, &c.DefaultValue, &c.PrimaryKey)
		if err != nil {
			return nil, err
		}

		columns = append(columns, c)
	}

	return columns, nil
}
