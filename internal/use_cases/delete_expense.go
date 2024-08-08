package usecases

import "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"

type DeleteExpenseInputDto struct {
	ExpenseID string `json:"id"`
}

type DeleteExpenseOutputDto struct {
	ID string `json:"id"`
}

type DeleteExpenseUseCase struct {
	ExpenseRepository repositories.ExpenseRepositoryInterface
}

func (c *DeleteExpenseUseCase) Execute(input DeleteExpenseInputDto) (DeleteExpenseOutputDto, []error) {
	expenseToDelete, err := c.ExpenseRepository.GetExpense(input.ExpenseID)
	if err != nil {
		return DeleteExpenseOutputDto{}, err
	}

	expenseToDelete.Deactivate()

	err = c.ExpenseRepository.DeleteExpense(expenseToDelete)
	if err != nil {
		return DeleteExpenseOutputDto{}, err
	}

	return DeleteExpenseOutputDto{
		ID: expenseToDelete.ID,
	}, nil
}
