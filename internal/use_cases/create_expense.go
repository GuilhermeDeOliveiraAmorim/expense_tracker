package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
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

func (c *CreateExpenseUseCase) Execute(input CreateExpenseInputDto) (CreateExpenseOutputDto, []util.ProblemDetails) {
	newExpense, err := entities.NewExpense(input.UserID, input.Amount, input.Category, input.Notes)
	if err != nil {
		return CreateExpenseOutputDto{}, err
	}

	CreateExpenseErr := c.ExpenseRepository.CreateExpense(*newExpense)
	if CreateExpenseErr != nil {
		return CreateExpenseOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Error creating new expense",
				Status:   500,
				Detail:   CreateExpenseErr.Error(),
				Instance: util.RFC500,
			},
		}
	}

	return CreateExpenseOutputDto{
		ID: newExpense.ID,
	}, nil
}
