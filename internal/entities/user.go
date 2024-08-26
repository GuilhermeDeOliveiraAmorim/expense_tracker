package entities

import (
	"time"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type User struct {
	SharedEntity
	Name  string `json:"name"`
	Login Login  `json:"login"`
}

func NewUser(name string, login Login) (*User, []util.ProblemDetails) {
	validationErrors := ValidateUser(name)

	if len(validationErrors) > 0 {
		return nil, validationErrors
	}

	return &User{
		SharedEntity: *NewSharedEntity(),
		Name:         name,
		Login:        login,
	}, nil
}

func ValidateUser(name string) []util.ProblemDetails {
	var validationErrors []util.ProblemDetails

	if name == "" {
		validationErrors = append(validationErrors, util.ProblemDetails{
			Type:     "Validation Error",
			Title:    "User name cannot be empty",
			Status:   400,
			Detail:   "User name is required",
			Instance: util.RFC400,
		})
	}

	if len(name) > 100 {
		validationErrors = append(validationErrors, util.ProblemDetails{
			Type:     "Validation Error",
			Title:    "User name is too long",
			Status:   400,
			Detail:   "User name must not exceed 100 characters",
			Instance: util.RFC400,
		})
	}

	return validationErrors
}

func (u *User) ChangeName(newName string) []util.ProblemDetails {
	validationErrors := ValidateUser(newName)

	if len(validationErrors) > 0 {
		return validationErrors
	}

	u.UpdatedAt = time.Now()
	u.Name = newName

	return validationErrors
}
