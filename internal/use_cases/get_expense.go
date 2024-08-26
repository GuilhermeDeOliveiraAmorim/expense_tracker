package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
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

func NewGetExpenseUseCase(
	ExpenseRepository repositories.ExpenseRepositoryInterface,
) *GetExpenseUseCase {
	return &GetExpenseUseCase{
		ExpenseRepository: ExpenseRepository,
	}
}

func (c *GetExpenseUseCase) Execute(input GetExpenseInputDto) (GetExpenseOutputDto, []util.ProblemDetails) {
	searchedExpense, err := c.ExpenseRepository.GetExpense(input.ExpenseID)
	if err != nil {
		return GetExpenseOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "Expense not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
			},
		}
	}

	return GetExpenseOutputDto{
		Expense: searchedExpense,
	}, nil
}
