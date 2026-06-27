package payload

type CreateFieldMetadataRequest struct {
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Required bool   `json:"required"`
	Nullable bool   `json:"nullable"`
	Order    uint   `json:"order"`
}

type CreateEntityMetadataRequest struct {
	Name           string                       `json:"name" binding:"required"`
	FieldsMetadata []CreateFieldMetadataRequest `json:"fields"`
}
