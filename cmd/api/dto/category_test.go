package dto

import (
	"testing"
)

func TestCategoryDTO_Validate(t *testing.T) {
	// Positive test case
	category := &CategoryDTO{
		Name: "Category Name",
	}
	err := category.Validate()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Negative test case
	emptyCategory := &CategoryDTO{}
	err = emptyCategory.Validate()
	if err == nil {
		t.Error("Expected an error, but got none")
	}
}
