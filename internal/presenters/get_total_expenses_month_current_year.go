package presenters

import (
	"strconv"
	"time"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type GetTotalExpensesMonthCurrentYearInputDto struct {
	UserID string `json:"user_id"`
	Year   string `json:"year"`
}

type GetTotalExpensesMonthCurrentYearOutputDto struct {
	ExpensesMonthCurrentYear repositories.ExpensesMonthCurrentYear `json:"expenses_month_current_year"`
}

type GetTotalExpensesMonthCurrentYearUseCase struct {
	PresentersRepository repositories.PresentersRepositoryInterface
	UserRepository       repositories.UserRepositoryInterface
}

func NewGetTotalExpensesMonthCurrentYearUseCase(
	PresentersRepository repositories.PresentersRepositoryInterface,
	UserRepository repositories.UserRepositoryInterface,
) *GetTotalExpensesMonthCurrentYearUseCase {
	return &GetTotalExpensesMonthCurrentYearUseCase{
		PresentersRepository: PresentersRepository,
		UserRepository:       UserRepository,
	}
}

func (c *GetTotalExpensesMonthCurrentYearUseCase) Execute(input GetTotalExpensesMonthCurrentYearInputDto) (GetTotalExpensesMonthCurrentYearOutputDto, []util.ProblemDetails) {
	user, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return GetTotalExpensesMonthCurrentYearOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
			},
		}
	} else if !user.Active {
		return GetTotalExpensesMonthCurrentYearOutputDto{}, []util.ProblemDetails{
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
		return GetTotalExpensesMonthCurrentYearOutputDto{}, []util.ProblemDetails{
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
		return GetTotalExpensesMonthCurrentYearOutputDto{}, []util.ProblemDetails{
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
		return GetTotalExpensesMonthCurrentYearOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Bad Request",
				Title:    "Invalid year",
				Status:   400,
				Detail:   "Year must be less than or equal to the current year",
				Instance: util.RFC400,
			},
		}
	}

	expensesMonthCurrentYear, getTotalExpensesMonthCurrentYearErr := c.PresentersRepository.GetTotalExpensesMonthCurrentYear(input.UserID, year)
	if getTotalExpensesMonthCurrentYearErr != nil {
		return GetTotalExpensesMonthCurrentYearOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Could not calculate total expenses",
				Status:   500,
				Detail:   getTotalExpensesMonthCurrentYearErr.Error(),
				Instance: util.RFC500,
			},
		}
	}

	return GetTotalExpensesMonthCurrentYearOutputDto{
		ExpensesMonthCurrentYear: expensesMonthCurrentYear,
	}, nil
}
