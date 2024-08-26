package entities

import (
	"time"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type Category struct {
	SharedEntity
	Name string `json:"name"`
}

func NewCategory(name string) (*Category, []util.ProblemDetails) {
	validationErrors := ValidateCategory(name)

	if len(validationErrors) > 0 {
		return nil, validationErrors
	}

	return &Category{
		SharedEntity: *NewSharedEntity(),
		Name:         name,
	}, nil
}

func ValidateCategory(name string) []util.ProblemDetails {
	var validationErrors []util.ProblemDetails

	if name == "" {
		validationErrors = append(validationErrors, util.ProblemDetails{
			Type:     "Validation Error",
			Title:    "Bad Request",
			Status:   400,
			Detail:   "Missing category name",
			Instance: util.RFC400,
		})
	}

	if len(name) > 100 {
		validationErrors = append(validationErrors, util.ProblemDetails{
			Type:     "Validation Error",
			Title:    "Bad Request",
			Status:   400,
			Detail:   "Category name cannot exceed 100 characters",
			Instance: util.RFC400,
		})
	}

	return validationErrors
}

func (c *Category) ChangeName(newName string) []util.ProblemDetails {
	validationErrors := ValidateCategory(newName)

	if len(validationErrors) > 0 {
		return validationErrors
	}

	c.UpdatedAt = time.Now()
	c.Name = newName

	return validationErrors
}
