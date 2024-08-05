package entities

import (
	"fmt"
	"time"
)

type Category struct {
	SharedEntity
	Name string `json:"name"`
}

func NewCategory(name string) (*Category, []error) {
	validationErrors := ValidateCategory(name)

	if len(validationErrors) > 0 {
		return nil, validationErrors
	}

	return &Category{
		SharedEntity: *NewSharedEntity(),
		Name:         name,
	}, nil
}

func ValidateCategory(name string) []error {
	var validationErrors []error

	if name == "" {
		validationErrors = append(validationErrors, fmt.Errorf("missing category name"))
	}

	if len(name) > 100 {
		validationErrors = append(validationErrors, fmt.Errorf("category name cannot exceed 100 characters"))
	}

	return validationErrors
}

func (c *Category) ChangeName(newName string) []error {
	validationErrors := ValidateCategory(newName)

	if len(validationErrors) > 0 {
		return validationErrors
	}

	c.UpdatedAt = time.Now()
	c.Name = newName

	return validationErrors
}
