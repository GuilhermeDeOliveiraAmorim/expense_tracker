package presenters

import (
	"strconv"
	"time"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type GetMonthlyExpensesByTagYearInputDto struct {
	UserID string `json:"user_id"`
	Year   string `json:"year"`
}

type GetMonthlyExpensesByTagYearOutputDto struct {
	Expenses       []repositories.MonthlyTagExpense `json:"expenses"`
	AvailableYears []int                            `json:"available_years"`
}

type GetMonthlyExpensesByTagYearUseCase struct {
	PresentersRepository repositories.PresentersRepositoryInterface
	UserRepository       repositories.UserRepositoryInterface
}

func NewGetMonthlyExpensesByTagYearUseCase(
	PresentersRepository repositories.PresentersRepositoryInterface,
	UserRepository repositories.UserRepositoryInterface,
) *GetMonthlyExpensesByTagYearUseCase {
	return &GetMonthlyExpensesByTagYearUseCase{
		PresentersRepository: PresentersRepository,
		UserRepository:       UserRepository,
	}
}

func (c *GetMonthlyExpensesByTagYearUseCase) Execute(input GetMonthlyExpensesByTagYearInputDto) (GetMonthlyExpensesByTagYearOutputDto, []util.ProblemDetails) {
	user, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return GetMonthlyExpensesByTagYearOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
			},
		}
	} else if !user.Active {
		return GetMonthlyExpensesByTagYearOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Forbidden",
				Title:    "User is not active",
				Status:   403,
				Detail:   "User is not active",
				Instance: util.RFC403,
			},
		}
	}

	year, errYear := strconv.Atoi(input.Year)
	if errYear != nil {
		return GetMonthlyExpensesByTagYearOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Bad Request",
				Title:    "Invalid year",
				Status:   400,
				Detail:   errYear.Error(),
				Instance: util.RFC400,
			},
		}
	}

	if year < 1900 || year > 99999 {
		return GetMonthlyExpensesByTagYearOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Bad Request",
				Title:    "Invalid year",
				Status:   400,
				Detail:   "Year must be between 1900 and 9999",
				Instance: util.RFC400,
			},
		}
	}

	if year > time.Now().Year() {
		return GetMonthlyExpensesByTagYearOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Bad Request",
				Title:    "Invalid year",
				Status:   400,
				Detail:   "Year must be less than or equal to the current year",
				Instance: util.RFC400,
			},
		}
	}

	expenses, availableYears, getMonthlyExpensesByTagYearErr := c.PresentersRepository.GetMonthlyExpensesByTagYear(input.UserID, year)
	if getMonthlyExpensesByTagYearErr != nil {
		return GetMonthlyExpensesByTagYearOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Could not calculate total expenses",
				Status:   500,
				Detail:   getMonthlyExpensesByTagYearErr.Error(),
				Instance: util.RFC500,
			},
		}
	}

	return GetMonthlyExpensesByTagYearOutputDto{
		Expenses:       expenses,
		AvailableYears: availableYears,
	}, nil
}
