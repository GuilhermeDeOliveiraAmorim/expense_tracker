package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
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

func (c *UpdateExpenseUseCase) Execute(input UpdateExpenseInputDto) (UpdateExpenseOutputDto, []error) {
	_, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return UpdateExpenseOutputDto{}, err
	}

	searchedExpense, err := c.ExpenseRepository.GetExpense(input.ExpenseID)
	if err != nil {
		return UpdateExpenseOutputDto{}, err
	}

	if input.Amount > 0 {
		err = searchedExpense.ChangeAmount(input.Amount)
		if len(err) > 0 {
			return UpdateExpenseOutputDto{}, err
		}
	}

	if input.Category.ID != "" {
		err = searchedExpense.ChangeCategory(input.Category)
		if len(err) > 0 {
			return UpdateExpenseOutputDto{}, err
		}
	}

	if input.Notes != "" {
		err = searchedExpense.ChangeNotes(input.Notes)
		if len(err) > 0 {
			return UpdateExpenseOutputDto{}, err
		}
	}

	err = c.ExpenseRepository.UpdateExpense(searchedExpense)
	if err != nil {
		return UpdateExpenseOutputDto{}, err
	}

	return UpdateExpenseOutputDto{
		ID: searchedExpense.ID,
	}, nil
}
