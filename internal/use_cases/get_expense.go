package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
)

type GetExpenseInputDto struct {
	ExpenseID string `json:"expense_id"`
}

type GetExpenseOutputDto struct {
	Expense entities.Expense `json:"expense"`
}

type GetExpenseUseCase struct {
	ExpenseRepository repositories.ExpenseRepositoryInterface
}

func (c *GetExpenseUseCase) Execute(input GetExpenseInputDto) (GetExpenseOutputDto, []error) {
	searchedExpense, err := c.ExpenseRepository.GetExpense(input.ExpenseID)
	if err != nil {
		return GetExpenseOutputDto{}, err
	}

	return GetExpenseOutputDto{
		Expense: searchedExpense,
	}, nil
}
