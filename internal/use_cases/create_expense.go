package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
)

type CreateExpenseInputDto struct {
	UserID   string            `json:"user_id"`
	Amount   float64           `json:"amount"`
	Category entities.Category `json:"category"`
	Notes    string            `json:"notes"`
}

type CreateExpenseOutputDto struct {
	ID string `json:"id"`
}

type CreateExpenseUseCase struct {
	ExpenseRepository repositories.ExpenseRepositoryInterface
}

func NewCreateExpenseUseCase(
	ExpenseRepository repositories.ExpenseRepositoryInterface,
) *CreateExpenseUseCase {
	return &CreateExpenseUseCase{
		ExpenseRepository: ExpenseRepository,
	}
}

func (c *CreateExpenseUseCase) Execute(input CreateExpenseInputDto) (CreateExpenseOutputDto, []error) {
	newExpense, err := entities.NewExpense(input.UserID, input.Amount, input.Category, input.Notes)
	if err != nil {
		return CreateExpenseOutputDto{}, err
	}

	err = c.ExpenseRepository.CreateExpense(*newExpense)
	if err != nil {
		return CreateExpenseOutputDto{}, err
	}

	return CreateExpenseOutputDto{
		ID: newExpense.ID,
	}, nil
}
