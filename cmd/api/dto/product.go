package dto

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type ProductDTO struct {
	ID          int64   `json:"id,omitempty"`
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"gte=0"`
	Image       string  `json:"image" validate:"required"`
	Category    string  `json:"category" validate:"required"`
}

type ProductUpdateDTO struct {
	Title        string  `json:"title" validate:"required"`
	Description  string  `json:"description" validate:"required"`
	Price        float64 `json:"price" validate:"gte=0"`
	Image        string  `json:"image" validate:"required"`
	CategoryName string  `json:"category" validate:"required"`
}

type ProductResponse struct {
	Data     []ProductDTO `json:"data"`
	Metadata Metadata     `json:"metadata"`
}

func (p *ProductDTO) Validate() error {
	return validate.Struct(p)
}
