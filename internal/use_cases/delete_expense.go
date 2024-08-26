package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type DeleteExpenseInputDto struct {
	ExpenseID string `json:"expense_id"`
}

type DeleteExpenseOutputDto struct {
	Message string `json:"message"`
}

type DeleteExpenseUseCase struct {
	ExpenseRepository repositories.ExpenseRepositoryInterface
}

func NewDeleteExpenseUseCase(
	ExpenseRepository repositories.ExpenseRepositoryInterface,
) *DeleteExpenseUseCase {
	return &DeleteExpenseUseCase{
		ExpenseRepository: ExpenseRepository,
	}
}

func (c *DeleteExpenseUseCase) Execute(input DeleteExpenseInputDto) (DeleteExpenseOutputDto, []util.ProblemDetails) {
	expenseToDelete, err := c.ExpenseRepository.GetExpense(input.ExpenseID)
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
		Message: "Expense deleted successfully",
	}, nil
}
