package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dotping-me/baastard/schema"
)

// Creates a lookup of tables by their table name
func columnMap(columns []schema.ColumnMetadata) map[string]schema.ColumnMetadata {

	m := make(map[string]schema.ColumnMetadata)
	for _, c := range columns {
		m[c.Name] = c
	}

	return m
}

// Helper that validates a given value type based on its SQL type representation.
func ValidateValue(sqlType string, value any) error {
	sqlType = strings.ToUpper(sqlType)

	switch sqlType {
	case "TEXT":
		if _, ok := value.(string); !ok {
			return errors.New("expected TEXT")
		}

	case "INTEGER":

		switch value.(type) {
		case float64:
		case int:
		case int32:
		case int64:

		default:
			return errors.New("expected INTEGER")
		}

	case "REAL":

		switch value.(type) {
		case float64:
		case float32:

		default:
			return errors.New("expected REAL")
		}

	case "BOOLEAN", "BOOL":
		if _, ok := value.(bool); !ok {
			return errors.New("expected BOOLEAN")
		}

	case "BLOB":
		// Skip validation for now.

	case "JSON":
		// map[string]any
		// []any
		// string

	default:
		// SQLite supports arbitrary type names.
		// Unknown types are accepted.
	}

	return nil
}

func ValidateInsert(columns []schema.ColumnMetadata, values map[string]any) error {

	m := columnMap(columns)
	for key := range values { // Checks whether field actually exists in table
		if _, ok := m[key]; !ok {
			return fmt.Errorf("unknown field: %s", key)
		}
	}

	// Validating required fields
	for _, col := range columns {
		if col.PrimaryKey {
			continue
		}

		if !col.NotNull {
			continue
		}

		if _, ok := values[col.Name]; !ok {
			return fmt.Errorf("missing required field: %s", col.Name)
		}
	}

	// Validating value types of each field
	for key, value := range values {
		err := ValidateValue(m[key].Type, value)
		if err != nil {
			return fmt.Errorf("%s: %w", key, err)
		}
	}

	return nil
}

func ValidateUpdate(columns []schema.ColumnMetadata, values map[string]any) error {

	m := columnMap(columns)
	for key, value := range values {
		col, ok := m[key]
		if !ok {
			return fmt.Errorf("unknown field: %s", key)
		}

		err := ValidateValue(col.Type, value)
		if err != nil {
			return fmt.Errorf("%s: %w", key, err)
		}
	}

	return nil
}
