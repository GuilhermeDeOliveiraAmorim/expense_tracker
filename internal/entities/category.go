package entities

import (
	"time"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type Category struct {
	SharedEntity
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
}

func NewCategory(userID string, name string, color string) (*Category, []util.ProblemDetails) {
	validationErrors := ValidateCategory(userID, name, color)

	if len(validationErrors) > 0 {
		return nil, validationErrors
	}

	return &Category{
		SharedEntity: *NewSharedEntity(),
		UserID:       userID,
		Name:         name,
		Color:        color,
	}, nil
}

func ValidateCategory(userID string, name string, color string) []util.ProblemDetails {
	var validationErrors []util.ProblemDetails

	if userID == "" {
		validationErrors = append(validationErrors, util.ProblemDetails{
			Type:     "Validation Error",
			Title:    "Bad Request",
			Status:   400,
			Detail:   "Missing user id",
			Instance: util.RFC400,
		})
	}

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

	if !util.IsValidHexColor(color) {
		validationErrors = append(validationErrors, util.ProblemDetails{
			Type:     "Validation Error",
			Title:    "Bad Request",
			Status:   400,
			Detail:   "Invalid color format. Use hexadecimal code (e.g., #FFFFFF or #FF0000)",
			Instance: util.RFC400,
		})
	}

	return validationErrors
}

func (c *Category) ChangeName(newName string) []util.ProblemDetails {
	var validationErrors []util.ProblemDetails

	if newName == "" {
		validationErrors = append(validationErrors, util.ProblemDetails{
			Type:     "Validation Error",
			Title:    "Bad Request",
			Status:   400,
			Detail:   "Missing category name",
			Instance: util.RFC400,
		})
	}

	if len(newName) > 100 {
		validationErrors = append(validationErrors, util.ProblemDetails{
			Type:     "Validation Error",
			Title:    "Bad Request",
			Status:   400,
			Detail:   "Category name cannot exceed 100 characters",
			Instance: util.RFC400,
		})
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	c.UpdatedAt = time.Now()
	c.Name = newName

	return validationErrors
}

func (c *Category) ChangeColor(newColor string) []util.ProblemDetails {
	var validationErrors []util.ProblemDetails

	if !util.IsValidHexColor(newColor) {
		validationErrors = append(validationErrors, util.ProblemDetails{
			Type:     "Validation Error",
			Title:    "Bad Request",
			Status:   400,
			Detail:   "Invalid color format. Use hexadecimal code (e.g., #FFFFFF or #FF0000)",
			Instance: util.RFC400,
		})
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	c.UpdatedAt = time.Now()
	c.Color = newColor

	return validationErrors
}
