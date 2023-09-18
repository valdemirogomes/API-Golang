package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductDTO_Validate(t *testing.T) {
	tests := []struct {
		name    string
		product *ProductDTO
		wantErr bool
	}{
		{
			name: "Test 1: Valid ProductDTO",
			product: &ProductDTO{
				Title:       "Test Product",
				Description: "Test Description",
				Price:       10.00,
				Image:       "test.jpg",
				Category:    "Test Category",
			},
			wantErr: false,
		},
		{
			name: "Test 2: Invalid ProductDTO",
			product: &ProductDTO{
				Description: "Test Description",
				Price:       10.00,
				Image:       "test.jpg",
				Category:    "Test Category",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.product.Validate()
			if tt.wantErr {
				assert.Error(t, err, "Error should be returned")
			} else {
				assert.NoError(t, err, "Error should not be returned")
			}

		})
	}
}
