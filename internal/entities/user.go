package entities

import (
	"fmt"
	"time"
)

type User struct {
	SharedEntity
	Name  string `json:"name"`
	Login Login  `json:"login"`
}

func NewUser(name string, login Login) (*User, []error) {
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

func ValidateUser(name string) []error {
	var validationErrors []error

	if name == "" {
		validationErrors = append(validationErrors, fmt.Errorf("missing user name"))
	}

	if len(name) > 100 {
		validationErrors = append(validationErrors, fmt.Errorf("user name cannot exceed 100 characters"))
	}

	return validationErrors
}

func (u *User) ChangeName(newName string) []error {
	validationErrors := ValidateUser(newName)

	if len(validationErrors) > 0 {
		return validationErrors
	}

	u.UpdatedAt = time.Now()
	u.Name = newName

	return validationErrors
}
