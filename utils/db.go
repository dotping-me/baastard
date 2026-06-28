package utils

import (
	"database/sql"
	"fmt"
	"strings"
)

// Universal helper that converts an arbitrary SQL result set into []map[string]any
func ScanRows(rows *sql.Rows) ([]map[string]any, error) {
	defer rows.Close()

	columns, err := rows.Columns() // Gets column names
	if err != nil {
		return nil, err
	}

	results := make([]map[string]any, 0)
	for rows.Next() {
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns)) // Scan requires pointers so yeah

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Reads one row.
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return nil, err
		}

		// Converts returned values to JSON-like structure
		row := make(map[string]any)
		for i, column := range columns {
			value := values[i]

			// SQLite returns TEXT as []byte.
			if bytes, ok := value.([]byte); ok {
				row[column] = string(bytes)
				continue
			}

			row[column] = value
		}

		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// Helper functions that just reduces redundant SQL command string building and
// generic SQL queries:
// 		SELECT, INSERT, UPDATE & DELETE

// TODO: Add validation for when building SQL commands

func Select(table string) string {
	return fmt.Sprintf("SELECT * FROM %s", table)
}

func SelectByID(table string) string {
	return fmt.Sprintf("SELECT * FROM %s WHERE id=? LIMIT 1", table)
}

func Delete(table string) string {
	return fmt.Sprintf("DELETE FROM %s WHERE id=?", table)
}

func Insert(table string, data map[string]any) (string, []any) {

	columns := make([]string, 0, len(data))
	values := make([]any, 0, len(data))
	placeholders := make([]string, 0, len(data))

	for column, value := range data {
		columns = append(columns, column)
		values = append(values, value)
		placeholders = append(placeholders, "?") // For ? (Position indices) in SQL Prepared Statement
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table,
		strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	return query, values
}

func Update(table string, id any, data map[string]any) (string, []any) {
	assignments := make([]string, 0, len(data))
	args := make([]any, 0, len(data)+1)

	for column, value := range data {
		assignments = append(assignments, column+"=?")
		args = append(args, value)
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=?", table,
		strings.Join(assignments, ", "))

	return query, args
}

func Exists() string {
	return "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?"
}

func Columns(table string) string {
	return fmt.Sprintf("PRAGMA table_info(%s)", table)
}
