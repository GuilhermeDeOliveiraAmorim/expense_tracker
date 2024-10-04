package presenters

import (
	"time"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type GetExpensesByMonthYearInputDto struct {
	UserID string `json:"user_id"`
	Month  int    `json:"month"`
	Year   int    `json:"year"`
}

type GetExpensesByMonthYearOutputDto struct {
	Expenses repositories.MonthExpenses `json:"expenses"`
}

type GetExpensesByMonthYearUseCase struct {
	PresentersRepository repositories.PresentersRepositoryInterface
	UserRepository       repositories.UserRepositoryInterface
}

func NewGetExpensesByMonthYearUseCase(
	PresentersRepository repositories.PresentersRepositoryInterface,
	UserRepository repositories.UserRepositoryInterface,
) *GetExpensesByMonthYearUseCase {
	return &GetExpensesByMonthYearUseCase{
		PresentersRepository: PresentersRepository,
		UserRepository:       UserRepository,
	}
}

func (c *GetExpensesByMonthYearUseCase) Execute(input GetExpensesByMonthYearInputDto) (GetExpensesByMonthYearOutputDto, []util.ProblemDetails) {
	user, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return GetExpensesByMonthYearOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
			},
		}
	} else if !user.Active {
		return GetExpensesByMonthYearOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Forbidden",
				Title:    "User is not active",
				Status:   403,
				Detail:   "User is not active",
				Instance: util.RFC403,
			},
		}
	}

	if input.Year < 1900 || input.Year > 99999 {
		return GetExpensesByMonthYearOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Bad Request",
				Title:    "Invalid year",
				Status:   400,
				Detail:   "Year must be between 1900 and 9999",
				Instance: util.RFC400,
			},
		}
	}

	if input.Year > time.Now().Year() {
		return GetExpensesByMonthYearOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Bad Request",
				Title:    "Invalid year",
				Status:   400,
				Detail:   "Year must be less than or equal to the current year",
				Instance: util.RFC400,
			},
		}
	}

	expenses, getExpensesByMonthYearErr := c.PresentersRepository.GetExpensesByMonthYear(input.UserID, input.Month, input.Year)
	if getExpensesByMonthYearErr != nil {
		return GetExpensesByMonthYearOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Could not calculate total expenses",
				Status:   500,
				Detail:   getExpensesByMonthYearErr.Error(),
				Instance: util.RFC500,
			},
		}
	}

	return GetExpensesByMonthYearOutputDto{
		Expenses: expenses,
	}, nil
}
