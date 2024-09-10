package usecases

import (
	"time"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type CreateExpenseInputDto struct {
	UserID      string  `json:"user_id"`
	Amount      float64 `json:"amount,string"`
	ExpenseDate string  `json:"expense_date"`
	CategoryID  string  `json:"category_id"`
	Notes       string  `json:"notes"`
}

type CreateExpenseOutputDto struct {
	ExpenseID string `json:"expense_id"`
}

type CreateExpenseUseCase struct {
	ExpenseRepository repositories.ExpenseRepositoryInterface
	UserRepository    repositories.UserRepositoryInterface
}

func NewCreateExpenseUseCase(
	ExpenseRepository repositories.ExpenseRepositoryInterface,
	UserRepository repositories.UserRepositoryInterface,
) *CreateExpenseUseCase {
	return &CreateExpenseUseCase{
		ExpenseRepository: ExpenseRepository,
		UserRepository:    UserRepository,
	}
}

func (c *CreateExpenseUseCase) Execute(input CreateExpenseInputDto) (CreateExpenseOutputDto, []util.ProblemDetails) {
	user, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return CreateExpenseOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
			},
		}
	} else if !user.Active {
		return CreateExpenseOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Forbidden",
				Title:    "User is not active",
				Status:   403,
				Detail:   "User is not active",
				Instance: util.RFC403,
			},
		}
	}

	location, err := time.LoadLocation(util.TIMEZONE)
	if err != nil {
		return CreateExpenseOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Validation Error",
				Title:    "Bad Request",
				Status:   400,
				Detail:   "Invalid timezone",
				Instance: util.RFC400,
			},
		}
	}

	newExpenseDate, err := time.ParseInLocation(util.DATEFORMAT, input.ExpenseDate, location)
	if err != nil {
		return CreateExpenseOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Validation Error",
				Title:    "Bad Request",
				Status:   400,
				Detail:   "Invalid expense date format",
				Instance: util.RFC400,
			},
		}
	}

	newExpense, NewExpenseErr := entities.NewExpense(input.UserID, input.Amount, newExpenseDate, input.CategoryID, input.Notes)
	if NewExpenseErr != nil {
		return CreateExpenseOutputDto{}, NewExpenseErr
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
		ExpenseID: newExpense.ID,
	}, nil
}
