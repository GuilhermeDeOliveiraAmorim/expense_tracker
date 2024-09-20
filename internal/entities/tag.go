package entities

import (
	"time"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type Tag struct {
	SharedEntity
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
}

func NewTag(userID string, name string, color string) (*Tag, []util.ProblemDetails) {
	validationErrors := ValidateTag(userID, name, color)

	if len(validationErrors) > 0 {
		return nil, validationErrors
	}

	return &Tag{
		SharedEntity: *NewSharedEntity(),
		UserID:       userID,
		Name:         name,
		Color:        color,
	}, nil
}

func ValidateTag(userID string, name string, color string) []util.ProblemDetails {
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
			Detail:   "Missing tag name",
			Instance: util.RFC400,
		})
	}

	if len(name) > 100 {
		validationErrors = append(validationErrors, util.ProblemDetails{
			Type:     "Validation Error",
			Title:    "Bad Request",
			Status:   400,
			Detail:   "Tag name cannot exceed 100 characters",
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

func (c *Tag) ChangeName(newName string) []util.ProblemDetails {
	var validationErrors []util.ProblemDetails

	if newName == "" {
		validationErrors = append(validationErrors, util.ProblemDetails{
			Type:     "Validation Error",
			Title:    "Bad Request",
			Status:   400,
			Detail:   "Missing tag name",
			Instance: util.RFC400,
		})
	}

	if len(newName) > 100 {
		validationErrors = append(validationErrors, util.ProblemDetails{
			Type:     "Validation Error",
			Title:    "Bad Request",
			Status:   400,
			Detail:   "Tag name cannot exceed 100 characters",
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

func (c *Tag) ChangeColor(newColor string) []util.ProblemDetails {
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
