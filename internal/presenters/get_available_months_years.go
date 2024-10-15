package presenters

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type GetAvailableMonthsYearsInputDto struct {
	UserID string `json:"user_id"`
}

type GetAvailableMonthsYearsOutputDto struct {
	AvailableYears  []int                      `json:"available_years"`
	AvailableMonths []repositories.MonthOption `json:"available_months"`
}

type GetAvailableMonthsYearsUseCase struct {
	PresentersRepository repositories.PresentersRepositoryInterface
	UserRepository       repositories.UserRepositoryInterface
}

func NewGetAvailableMonthsYearsUseCase(
	PresentersRepository repositories.PresentersRepositoryInterface,
	UserRepository repositories.UserRepositoryInterface,
) *GetAvailableMonthsYearsUseCase {
	return &GetAvailableMonthsYearsUseCase{
		PresentersRepository: PresentersRepository,
		UserRepository:       UserRepository,
	}
}

func (c *GetAvailableMonthsYearsUseCase) Execute(input GetAvailableMonthsYearsInputDto) (GetAvailableMonthsYearsOutputDto, []util.ProblemDetails) {
	user, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return GetAvailableMonthsYearsOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
			},
		}
	} else if !user.Active {
		return GetAvailableMonthsYearsOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Forbidden",
				Title:    "User is not active",
				Status:   403,
				Detail:   "User is not active",
				Instance: util.RFC403,
			},
		}
	}

	availableYears, availableMonths, getExpensesByMonthYearErr := c.PresentersRepository.GetAvailableMonthsYears(input.UserID)
	if getExpensesByMonthYearErr != nil {
		return GetAvailableMonthsYearsOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Could not calculate total expenses",
				Status:   500,
				Detail:   getExpensesByMonthYearErr.Error(),
				Instance: util.RFC500,
			},
		}
	}

	return GetAvailableMonthsYearsOutputDto{
		AvailableYears:  availableYears,
		AvailableMonths: availableMonths,
	}, nil
}
