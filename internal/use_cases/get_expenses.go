package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type GetExpensesInputDto struct {
}

type GetExpensesOutputDto struct {
	Expenses []entities.Expense `json:"expenses"`
}

type GetExpensesUseCase struct {
	ExpenseRepository repositories.ExpenseRepositoryInterface
}

func NewGetExpensesUseCase(
	ExpenseRepository repositories.ExpenseRepositoryInterface,
) *GetExpensesUseCase {
	return &GetExpensesUseCase{
		ExpenseRepository: ExpenseRepository,
	}
}

func (c *GetExpensesUseCase) Execute(input GetExpensesInputDto) (GetExpensesOutputDto, []util.ProblemDetails) {
	searchedsExpenses, err := c.ExpenseRepository.GetExpenses()
	if err != nil {
		return GetExpensesOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "An error occurred while retrieving expenses",
				Status:   500,
				Detail:   err.Error(),
				Instance: util.RFC500,
			},
		}
	}

	return GetExpensesOutputDto{
		Expenses: searchedsExpenses,
	}, nil
}
