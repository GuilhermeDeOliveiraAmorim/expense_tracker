package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type UpdateExpenseInputDto struct {
	UserID    string            `json:"user_id"`
	ExpenseID string            `json:"id"`
	Amount    float64           `json:"amount"`
	Category  entities.Category `json:"category"`
	Notes     string            `json:"notes"`
}

type UpdateExpenseOutputDto struct {
	ID string `json:"id"`
}

type UpdateExpenseUseCase struct {
	ExpenseRepository repositories.ExpenseRepositoryInterface
	UserRepository    repositories.UserRepositoryInterface
}

func NewUpdateExpenseUseCase(
	ExpenseRepository repositories.ExpenseRepositoryInterface,
) *UpdateExpenseUseCase {
	return &UpdateExpenseUseCase{
		ExpenseRepository: ExpenseRepository,
	}
}

func (c *UpdateExpenseUseCase) Execute(input UpdateExpenseInputDto) (UpdateExpenseOutputDto, []util.ProblemDetails) {
	_, err := c.UserRepository.GetUser(input.UserID)
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
	}

	searchedExpense, GetExpenseErr := c.ExpenseRepository.GetExpense(input.ExpenseID)
	if GetExpenseErr != nil {
		return UpdateExpenseOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "Expense not found",
				Status:   404,
				Detail:   GetExpenseErr.Error(),
				Instance: util.RFC404,
			},
		}
	}

	if input.Amount > 0 {
		err := searchedExpense.ChangeAmount(input.Amount)
		if len(err) > 0 {
			return UpdateExpenseOutputDto{}, err
		}
	}

	if input.Category.ID != "" {
		err := searchedExpense.ChangeCategory(input.Category)
		if len(err) > 0 {
			return UpdateExpenseOutputDto{}, err
		}
	}

	if input.Notes != "" {
		err := searchedExpense.ChangeNotes(input.Notes)
		if len(err) > 0 {
			return UpdateExpenseOutputDto{}, err
		}
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
		ID: searchedExpense.ID,
	}, nil
}
