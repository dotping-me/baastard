package schema

// Used to set which table is exposed through API endpoints and which are not.

type TableConfig struct {
	Name          string `json:"name"          gorm:"primaryKey"`
	Hidden        bool   `json:"hidden"        gorm:"default:false"` // TODO: Modify so that only specific rows can be hidden and not the entire table
	ReadOnly      bool   `json:"readOnly"      gorm:"default:false"` // ReadOnly by a user
	AuthRequired  bool   `json:"authRequired"  gorm:"default:true"`
	VectorEnabled bool   `json:"vectorEnabled" gorm:"default:false"`

	// TODO: Add ability to specify which fields are encrypted
}

type ColumnMetadata struct {
	CID          int
	Name         string
	Type         string
	NotNull      bool    `gorm:"column:notnull"`
	DefaultValue *string `gorm:"column:dflt_value"`
	PrimaryKey   bool    `gorm:"column:pk"`
}
