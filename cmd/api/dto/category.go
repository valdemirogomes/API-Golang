package dto

type CategoryDTO struct {
	Name string `json:"name" validate:"required"`
}

type CategoryResponseDTO struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name"`
}

type CategoryListResponseDTO struct {
	Data     []CategoryResponseDTO `json:"data"`
	Metadata Metadata              `json:"metadata"`
}

func (c *CategoryDTO) Validate() error {
	return validate.Struct(c)
}
