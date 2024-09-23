package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type GetExpenseInputDto struct {
	UserID    string `json:"user_id"`
	ExpenseID string `json:"expense_id"`
}

type GetExpenseOutputDto struct {
	Expense entities.Expense `json:"expense"`
}

type GetExpenseUseCase struct {
	ExpenseRepository repositories.ExpenseRepositoryInterface
	UserRepository    repositories.UserRepositoryInterface
}

func NewGetExpenseUseCase(
	ExpenseRepository repositories.ExpenseRepositoryInterface,
	UserRepository repositories.UserRepositoryInterface,
) *GetExpenseUseCase {
	return &GetExpenseUseCase{
		ExpenseRepository: ExpenseRepository,
		UserRepository:    UserRepository,
	}
}

func (c *GetExpenseUseCase) Execute(input GetExpenseInputDto) (GetExpenseOutputDto, []util.ProblemDetails) {
	user, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return GetExpenseOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
			},
		}
	} else if !user.Active {
		return GetExpenseOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Forbidden",
				Title:    "User is not active",
				Status:   403,
				Detail:   "User is not active",
				Instance: util.RFC403,
			},
		}
	}

	searchedExpense, err := c.ExpenseRepository.GetExpense(input.UserID, input.ExpenseID)
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
