package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
)

type GetExpensesInputDto struct {
}

type GetExpensesOutputDto struct {
	Expenses []entities.Expense `json:"categories"`
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

func (c *GetExpensesUseCase) Execute(input GetExpensesInputDto) (GetExpensesOutputDto, []error) {
	searchedsExpenses, err := c.ExpenseRepository.GetExpenses()
	if err != nil {
		return GetExpensesOutputDto{}, err
	}

	return GetExpensesOutputDto{
		Expenses: searchedsExpenses,
	}, nil
}
