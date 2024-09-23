package entities

import (
	"time"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type Expense struct {
	SharedEntity
	UserID      string    `json:"user_id"`
	Amount      float64   `json:"amount,string,omitempty"`
	ExpenseDate time.Time `json:"expense_date"`
	CategoryID  string    `json:"category_id"`
	TagIDs      []string  `json:"tag_ids"`
	Notes       string    `json:"notes"`
	Category    Category  `json:"category"`
	Tags        []Tag     `json:"tags"`
}

func NewExpense(userID string, amount float64, expenseDate time.Time, categoryID string, notes string) (*Expense, []util.ProblemDetails) {
	validationErrors := ValidateExpense(userID, amount, categoryID, notes)

	if len(validationErrors) > 0 {
		return nil, validationErrors
	}

	return &Expense{
		SharedEntity: *NewSharedEntity(),
		UserID:       userID,
		Amount:       amount,
		ExpenseDate:  expenseDate,
		CategoryID:   categoryID,
		Notes:        notes,
	}, nil
}

func ValidateExpense(userID string, amount float64, categoryID string, notes string) []util.ProblemDetails {
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

	if categoryID == "" {
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

	if len(validationErrors) > 0 {
		return validationErrors
	} else {
		e.UpdatedAt = time.Now()
		e.Amount = newAmount

		return validationErrors
	}
}

func (e *Expense) AddTagByID(tagID string) []util.ProblemDetails {
	var validationErrors []util.ProblemDetails

	if tagID == "" {
		validationErrors = append(validationErrors, util.ProblemDetails{
			Type:     "Validation Error",
			Title:    "Invalid Tag ID",
			Status:   400,
			Detail:   "Tag ID cannot be empty",
			Instance: util.RFC400,
		})
		return validationErrors
	}

	e.TagIDs = append(e.TagIDs, tagID)

	return validationErrors
}

func (e *Expense) ChangeExpenseDate(expenseDate string) []util.ProblemDetails {
	var validationErrors []util.ProblemDetails

	newExpenseDate, err := time.Parse("02012006", expenseDate)
	if err != nil {
		validationErrors = append(validationErrors, util.ProblemDetails{
			Type:     "Validation Error",
			Title:    "Bad Request",
			Status:   400,
			Detail:   "Invalid expense date format",
			Instance: util.RFC400,
		})
	}

	if len(validationErrors) > 0 {
		return validationErrors
	} else {
		e.UpdatedAt = time.Now()
		e.ExpenseDate = newExpenseDate

		return validationErrors
	}
}

func (e *Expense) ChangeCategory(newCategoryID string) []util.ProblemDetails {
	var validationErrors []util.ProblemDetails

	if newCategoryID == "" {
		validationErrors = append(validationErrors, util.ProblemDetails{
			Type:     "Validation Error",
			Title:    "Bad Request",
			Status:   400,
			Detail:   "Missing new category ID",
			Instance: util.RFC400,
		})
	}

	if len(validationErrors) > 0 {
		return validationErrors
	} else {
		e.UpdatedAt = time.Now()
		e.CategoryID = newCategoryID

		return validationErrors
	}
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

	if len(validationErrors) > 0 {
		return validationErrors
	} else {
		e.UpdatedAt = time.Now()
		e.Notes = newNotes

		return validationErrors
	}
}
