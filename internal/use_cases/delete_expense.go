package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type DeleteExpenseInputDto struct {
	UserID    string `json:"user_id"`
	ExpenseID string `json:"expense_id"`
}

type DeleteExpenseOutputDto struct {
	SuccessMessage string `json:"success_message"`
	ContentMessage string `json:"content_message"`
}

type DeleteExpenseUseCase struct {
	ExpenseRepository repositories.ExpenseRepositoryInterface
	UserRepository    repositories.UserRepositoryInterface
}

func NewDeleteExpenseUseCase(
	ExpenseRepository repositories.ExpenseRepositoryInterface,
	UserRepository repositories.UserRepositoryInterface,
) *DeleteExpenseUseCase {
	return &DeleteExpenseUseCase{
		ExpenseRepository: ExpenseRepository,
		UserRepository:    UserRepository,
	}
}

func (c *DeleteExpenseUseCase) Execute(input DeleteExpenseInputDto) (DeleteExpenseOutputDto, []util.ProblemDetails) {
	user, getUserErr := c.UserRepository.GetUser(input.UserID)
	if getUserErr != nil {
		return DeleteExpenseOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   getUserErr.Error(),
				Instance: util.RFC404,
			},
		}
	} else if !user.Active {
		return DeleteExpenseOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Forbidden",
				Title:    "User is not active",
				Status:   403,
				Detail:   "User is not active",
				Instance: util.RFC403,
			},
		}
	}

	expenseToDelete, getExpenseErr := c.ExpenseRepository.GetExpense(input.UserID, input.ExpenseID)
	if getExpenseErr != nil {
		return DeleteExpenseOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "Expense not found",
				Status:   404,
				Detail:   getExpenseErr.Error(),
				Instance: util.RFC404,
			},
		}
	}

	expenseToDelete.Expense.Deactivate()

	deleteExpenseErr := c.ExpenseRepository.DeleteExpense(expenseToDelete.Expense)
	if deleteExpenseErr != nil {
		return DeleteExpenseOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Err deleting expense",
				Status:   500,
				Detail:   deleteExpenseErr.Error(),
				Instance: util.RFC500,
			},
		}
	}

	return DeleteExpenseOutputDto{
		SuccessMessage: "Expense deleted successfully",
		ContentMessage: "Expense ID: " + input.ExpenseID,
	}, nil
}
