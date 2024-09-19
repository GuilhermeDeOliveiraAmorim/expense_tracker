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
	user, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return DeleteExpenseOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   err.Error(),
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

	expenseToDelete, err := c.ExpenseRepository.GetExpense(input.UserID, input.ExpenseID)
	if err != nil {
		return DeleteExpenseOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "Expense not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
			},
		}
	}

	expenseToDelete.Deactivate()

	DeleteExpenseErr := c.ExpenseRepository.DeleteExpense(expenseToDelete)
	if DeleteExpenseErr != nil {
		return DeleteExpenseOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Err deleting expense",
				Status:   500,
				Detail:   DeleteExpenseErr.Error(),
				Instance: util.RFC500,
			},
		}
	}

	return DeleteExpenseOutputDto{
		SuccessMessage: "Expense deleted successfully",
		ContentMessage: "Expense ID: " + input.ExpenseID,
	}, nil
}
