package presenters

import (
	"time"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type GetCategoryTagsTotalsByMonthYearInputDto struct {
	UserID string `json:"user_id"`
	Month  int    `json:"month"`
	Year   int    `json:"year"`
}

type GetCategoryTagsTotalsByMonthYearOutputDto struct {
	Expenses repositories.CategoryTagsTotals `json:"expenses"`
}

type GetCategoryTagsTotalsByMonthYearUseCase struct {
	PresentersRepository repositories.PresentersRepositoryInterface
	UserRepository       repositories.UserRepositoryInterface
}

func NewGetCategoryTagsTotalsByMonthYearUseCase(
	PresentersRepository repositories.PresentersRepositoryInterface,
	UserRepository repositories.UserRepositoryInterface,
) *GetCategoryTagsTotalsByMonthYearUseCase {
	return &GetCategoryTagsTotalsByMonthYearUseCase{
		PresentersRepository: PresentersRepository,
		UserRepository:       UserRepository,
	}
}

func (c *GetCategoryTagsTotalsByMonthYearUseCase) Execute(input GetCategoryTagsTotalsByMonthYearInputDto) (GetCategoryTagsTotalsByMonthYearOutputDto, []util.ProblemDetails) {
	user, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return GetCategoryTagsTotalsByMonthYearOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
			},
		}
	} else if !user.Active {
		return GetCategoryTagsTotalsByMonthYearOutputDto{}, []util.ProblemDetails{
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
		return GetCategoryTagsTotalsByMonthYearOutputDto{}, []util.ProblemDetails{
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
		return GetCategoryTagsTotalsByMonthYearOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Bad Request",
				Title:    "Invalid year",
				Status:   400,
				Detail:   "Year must be less than or equal to the current year",
				Instance: util.RFC400,
			},
		}
	}

	expenses, getExpensesByMonthYearErr := c.PresentersRepository.GetCategoryTagsTotalsByMonthYear(input.UserID, input.Month, input.Year)
	if getExpensesByMonthYearErr != nil {
		return GetCategoryTagsTotalsByMonthYearOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Could not calculate total expenses",
				Status:   500,
				Detail:   getExpensesByMonthYearErr.Error(),
				Instance: util.RFC500,
			},
		}
	}

	return GetCategoryTagsTotalsByMonthYearOutputDto{
		Expenses: expenses,
	}, nil
}
