package presenters

import (
	"time"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type GetMonthlyExpensesByCategoryPeriodInputDto struct {
	UserID string `json:"user_id"`
	Year   int    `json:"year"`
}

type GetMonthlyExpensesByCategoryPeriodOutputDto struct {
	Expenses       []repositories.MonthlyCategoryExpense `json:"expenses"`
	AvailableYears []int                                 `json:"available_years"`
}

type GetMonthlyExpensesByCategoryPeriodUseCase struct {
	PresentersRepository repositories.PresentersRepositoryInterface
	UserRepository       repositories.UserRepositoryInterface
}

func NewGetMonthlyExpensesByCategoryPeriodUseCase(
	PresentersRepository repositories.PresentersRepositoryInterface,
	UserRepository repositories.UserRepositoryInterface,
) *GetMonthlyExpensesByCategoryPeriodUseCase {
	return &GetMonthlyExpensesByCategoryPeriodUseCase{
		PresentersRepository: PresentersRepository,
		UserRepository:       UserRepository,
	}
}

func (c *GetMonthlyExpensesByCategoryPeriodUseCase) Execute(input GetMonthlyExpensesByCategoryPeriodInputDto) (GetMonthlyExpensesByCategoryPeriodOutputDto, []util.ProblemDetails) {
	user, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return GetMonthlyExpensesByCategoryPeriodOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
			},
		}
	} else if !user.Active {
		return GetMonthlyExpensesByCategoryPeriodOutputDto{}, []util.ProblemDetails{
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
		return GetMonthlyExpensesByCategoryPeriodOutputDto{}, []util.ProblemDetails{
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
		return GetMonthlyExpensesByCategoryPeriodOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Bad Request",
				Title:    "Invalid year",
				Status:   400,
				Detail:   "Year must be less than or equal to the current year",
				Instance: util.RFC400,
			},
		}
	}

	expenses, availableYears, getMonthlyExpensesByCategoryPeriodErr := c.PresentersRepository.GetMonthlyExpensesByCategoryPeriod(input.UserID, input.Year)
	if getMonthlyExpensesByCategoryPeriodErr != nil {
		return GetMonthlyExpensesByCategoryPeriodOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Could not calculate total expenses",
				Status:   500,
				Detail:   getMonthlyExpensesByCategoryPeriodErr.Error(),
				Instance: util.RFC500,
			},
		}
	}

	return GetMonthlyExpensesByCategoryPeriodOutputDto{
		Expenses:       expenses,
		AvailableYears: availableYears,
	}, nil
}
