package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type GetExpensesInputDto struct {
	UserID string `json:"user_id"`
}

type GetExpensesOutputDto struct {
	Expenses []repositories.GetExpense `json:"expenses"`
}

type GetExpensesUseCase struct {
	ExpenseRepository repositories.ExpenseRepositoryInterface
	UserRepository    repositories.UserRepositoryInterface
}

func NewGetExpensesUseCase(
	ExpenseRepository repositories.ExpenseRepositoryInterface,
	UserRepository repositories.UserRepositoryInterface,
) *GetExpensesUseCase {
	return &GetExpensesUseCase{
		ExpenseRepository: ExpenseRepository,
		UserRepository:    UserRepository,
	}
}

func (c *GetExpensesUseCase) Execute(input GetExpensesInputDto) (GetExpensesOutputDto, []util.ProblemDetails) {
	user, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return GetExpensesOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
			},
		}
	} else if !user.Active {
		return GetExpensesOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Forbidden",
				Title:    "User is not active",
				Status:   403,
				Detail:   "User is not active",
				Instance: util.RFC403,
			},
		}
	}

	searchedsExpenses, err := c.ExpenseRepository.GetExpenses(input.UserID)
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
