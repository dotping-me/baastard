package services

import (
	"errors"

	"github.com/dotping-me/baastard/schema"
	"gorm.io/gorm"
)

// When a new table is created, this service syncs up table_configs table;
// autopopulating it.

type TableConfigService struct {
	db *gorm.DB
}

var ignoredTables = map[string]bool{
	"sqlite_sequence": true,
	"sqlite_schema":   true,
	"table_configs":   true,
	"users":           true,
	"sessions":        true,
}

func NewTableConfigService(db *gorm.DB) *TableConfigService {
	return &TableConfigService{db: db}
}

// Syncs up the table_configs table to match database state on startup because
// SQL Triggers do not run on CREATE/DROP. NOTE: This should be called at runtime
// too when a table is created or dropped.
func (s *TableConfigService) Sync() error {

	// Queries all table names in the database
	var tables []string
	err := s.db.Raw(`SELECT name FROM sqlite_master	
		WHERE type='table'`).Scan(&tables).Error
	if err != nil {
		return err
	}

	actualTables := make(map[string]bool)
	for _, table := range tables {
		if ignoredTables[table] {
			continue
		}

		actualTables[table] = true
	}

	// Finds which table has an existing configuration for it
	var configs []schema.TableConfig
	err = s.db.Find(&configs).Error
	if err != nil {
		return err
	}

	configuredTables := make(map[string]bool)
	for _, cfg := range configs {
		configuredTables[cfg.Name] = true
	}

	// Syncs up table_configs table
	for t := range actualTables {
		if configuredTables[t] {
			continue
		}

		config := schema.TableConfig{Name: t}
		err = s.db.Create(&config).Error
		if err != nil {
			return err
		}
	}

	// Cleans up state; some tables may have been deleted
	for _, config := range configs {
		if ignoredTables[config.Name] {
			continue
		}

		if actualTables[config.Name] {
			continue
		}

		err = s.db.Delete(&config).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// Returns the configuration for a table with a given name
func (s *TableConfigService) Get(name string) (*schema.TableConfig, error) {

	var cfg schema.TableConfig
	err := s.db.Where("name = ?", name).First(&cfg).Error
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Returns all table configurations
func (s *TableConfigService) GetAll() ([]schema.TableConfig, error) {

	var configs []schema.TableConfig
	err := s.db.Find(&configs).Error
	if err != nil {
		return nil, err
	}

	return configs, nil
}

// Updates a table configuration.
func (s *TableConfigService) Update(cfg *schema.TableConfig) error {
	return s.db.Save(cfg).Error
}

// Creating a configuration manually. Normally Sync() would handle this but it's
// good to still have it just in case
func (s *TableConfigService) Create(cfg *schema.TableConfig) error {

	var existing schema.TableConfig
	err := s.db.Where("name = ?", cfg.Name).First(&existing).Error
	if err == nil {
		return errors.New("table configuration already exists")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return s.db.Create(cfg).Error
}

// Deletes a configuration
func (s *TableConfigService) Delete(name string) error {
	return s.db.Where("name = ?", name).Delete(&schema.TableConfig{}).Error
}
