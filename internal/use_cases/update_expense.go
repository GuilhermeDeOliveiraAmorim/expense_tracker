package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type UpdateExpenseInputDto struct {
	UserID      string   `json:"user_id"`
	ExpenseID   string   `json:"expense_id"`
	Amount      float64  `json:"amount,string"`
	ExpenseDate string   `json:"expense_date"`
	CategoryID  string   `json:"category_id"`
	Notes       string   `json:"notes"`
	Tags        []string `json:"tags"`
}

type UpdateExpenseOutputDto struct {
	ExpenseID      string `json:"expense_id"`
	SuccessMessage string `json:"success_message"`
	ContentMessage string `json:"content_message"`
}

type UpdateExpenseUseCase struct {
	ExpenseRepository repositories.ExpenseRepositoryInterface
	UserRepository    repositories.UserRepositoryInterface
}

func NewUpdateExpenseUseCase(
	ExpenseRepository repositories.ExpenseRepositoryInterface,
	UserRepository repositories.UserRepositoryInterface,
) *UpdateExpenseUseCase {
	return &UpdateExpenseUseCase{
		ExpenseRepository: ExpenseRepository,
		UserRepository:    UserRepository,
	}
}

func (c *UpdateExpenseUseCase) Execute(input UpdateExpenseInputDto) (UpdateExpenseOutputDto, []util.ProblemDetails) {
	var validationErrors []util.ProblemDetails

	if input.UserID == "" {
		validationErrors = append(validationErrors, util.ProblemDetails{
			Type:     "Validation Error",
			Title:    "Invalid User ID",
			Status:   400,
			Detail:   "User ID cannot be empty",
			Instance: util.RFC400,
		})
	}
	if input.ExpenseID == "" {
		validationErrors = append(validationErrors, util.ProblemDetails{
			Type:     "Validation Error",
			Title:    "Invalid Expense ID",
			Status:   400,
			Detail:   "Expense ID cannot be empty",
			Instance: util.RFC400,
		})
	}

	if len(validationErrors) > 0 {
		return UpdateExpenseOutputDto{}, validationErrors
	}

	user, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return UpdateExpenseOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
			},
		}
	} else if !user.Active {
		return UpdateExpenseOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Forbidden",
				Title:    "User is not active",
				Status:   403,
				Detail:   "User is not active",
				Instance: util.RFC403,
			},
		}
	}

	searchedExpense, GetExpenseErr := c.ExpenseRepository.GetExpense(input.UserID, input.ExpenseID)
	if GetExpenseErr != nil {
		return UpdateExpenseOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "Expense or not found",
				Status:   404,
				Detail:   GetExpenseErr.Error(),
				Instance: util.RFC404,
			},
		}
	}

	if input.Amount > 0 {
		err := searchedExpense.ChangeAmount(input.Amount)
		if len(err) > 0 {
			validationErrors = append(validationErrors, err...)
		}
	}

	if input.ExpenseDate != "" {
		err := searchedExpense.ChangeExpenseDate(input.ExpenseDate)
		if len(err) > 0 {
			validationErrors = append(validationErrors, err...)
		}
	}

	if input.CategoryID != "" {
		err := searchedExpense.ChangeCategory(input.CategoryID)
		if len(err) > 0 {
			validationErrors = append(validationErrors, err...)
		}
	}

	changeNotesErr := searchedExpense.ChangeNotes(input.Notes)
	if len(changeNotesErr) > 0 {
		validationErrors = append(validationErrors, changeNotesErr...)
	}

	changeTagsErr := searchedExpense.ChangeTags(input.Tags)
	if len(changeTagsErr) > 0 {
		validationErrors = append(validationErrors, changeTagsErr...)
	}

	if len(validationErrors) > 0 {
		return UpdateExpenseOutputDto{}, validationErrors
	}

	UpdateExpenseErr := c.ExpenseRepository.UpdateExpense(searchedExpense)
	if UpdateExpenseErr != nil {
		return UpdateExpenseOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "An error occurred while updating expense",
				Status:   500,
				Detail:   UpdateExpenseErr.Error(),
				Instance: util.RFC500,
			},
		}
	}

	return UpdateExpenseOutputDto{
		ExpenseID:      input.ExpenseID,
		SuccessMessage: "Expense updated successfully",
		ContentMessage: "Expense ID: " + input.ExpenseID,
	}, nil
}
