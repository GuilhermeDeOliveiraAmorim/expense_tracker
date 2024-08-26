package entities

import (
	"time"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type Expense struct {
	SharedEntity
	UserID   string   `json:"user_id"`
	Amount   float64  `json:"amount"`
	Category Category `json:"category"`
	Notes    string   `json:"notes"`
}

func NewExpense(userID string, amount float64, category Category, notes string) (*Expense, []util.ProblemDetails) {
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

func ValidateExpense(userID string, amount float64, category Category, notes string) []util.ProblemDetails {
	var validationErrors []util.ProblemDetails

	if userID == "" {
		validationErrors = append(validationErrors, util.ProblemDetails{
			Type:     "Validation Error",
			Title:    "Bad Request",
			Status:   400,
			Detail:   "Missing user ID",
			Instance: util.RFC400,
		})
	}

	if amount <= 0 {
		validationErrors = append(validationErrors, util.ProblemDetails{
			Type:     "Validation Error",
			Title:    "Bad Request",
			Status:   400,
			Detail:   "Amount must be greater than 0",
			Instance: util.RFC400,
		})
	}

	if category.ID == "" {
		validationErrors = append(validationErrors, util.ProblemDetails{
			Type:     "Validation Error",
			Title:    "Bad Request",
			Status:   400,
			Detail:   "Missing category ID",
			Instance: util.RFC400,
		})
	}

	if len(notes) > 200 {
		validationErrors = append(validationErrors, util.ProblemDetails{
			Type:     "Validation Error",
			Title:    "Bad Request",
			Status:   400,
			Detail:   "Notes cannot exceed 200 characters",
			Instance: util.RFC400,
		})
	}

	return validationErrors
}

func (e *Expense) ChangeAmount(newAmount float64) []util.ProblemDetails {
	var validationErrors []util.ProblemDetails

	if newAmount <= 0 {
		validationErrors = append(validationErrors, util.ProblemDetails{
			Type:     "Validation Error",
			Title:    "Bad Request",
			Status:   400,
			Detail:   "New amount must be greater than 0",
			Instance: util.RFC400,
		})
	}

	e.UpdatedAt = time.Now()
	e.Amount = newAmount

	return validationErrors
}

func (e *Expense) ChangeCategory(newCategory Category) []util.ProblemDetails {
	var validationErrors []util.ProblemDetails

	if newCategory.ID == "" {
		validationErrors = append(validationErrors, util.ProblemDetails{
			Type:     "Validation Error",
			Title:    "Bad Request",
			Status:   400,
			Detail:   "Missing new category ID",
			Instance: util.RFC400,
		})
	}

	e.UpdatedAt = time.Now()
	e.Category = newCategory

	return validationErrors
}

func (e *Expense) ChangeNotes(newNotes string) []util.ProblemDetails {
	var validationErrors []util.ProblemDetails

	if len(newNotes) > 200 {
		validationErrors = append(validationErrors, util.ProblemDetails{
			Type:     "Validation Error",
			Title:    "Bad Request",
			Status:   400,
			Detail:   "New notes cannot exceed 200 characters",
			Instance: util.RFC400,
		})
	}

	e.UpdatedAt = time.Now()
	e.Notes = newNotes

	return validationErrors
}
