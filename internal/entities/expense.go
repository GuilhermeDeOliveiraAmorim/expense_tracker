package entities

import (
	"fmt"
)

type Expense struct {
	SharedEntity
	UserID   string   `json:"user_id"`
	Amount   float64  `json:"amount"`
	Category Category `json:"category"`
	Notes    string   `json:"notes"`
}

func NewExpense(userID string, amount float64, category Category, notes string) (*Expense, []error) {
	validationErrors := ValidateExpense(userID, amount, category, notes)

	if len(validationErrors) > 0 {
		return nil, validationErrors
	}

	return &Expense{
		SharedEntity: *NewSharedEntity(),
		UserID:       userID,
		Amount:       amount,
		Category:     category,
		Notes:        notes,
	}, nil
}

func ValidateExpense(userID string, amount float64, category Category, notes string) []error {
	var validationErrors []error

	if userID == "" {
		validationErrors = append(validationErrors, fmt.Errorf("missing user ID"))
	}

	if amount <= 0 {
		validationErrors = append(validationErrors, fmt.Errorf("amount must be greater than 0"))
	}

	if category.ID == "" {
		validationErrors = append(validationErrors, fmt.Errorf("missing category ID"))
	}

	if len(notes) > 200 {
		validationErrors = append(validationErrors, fmt.Errorf("notes cannot exceed 200 characters"))
	}

	return validationErrors
}
