package presenters

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type GetTotalExpensesForPeriodInputDto struct {
	UserID    string `json:"user_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type GetTotalExpensesForPeriodOutputDto struct {
	Total float64 `json:"total"`
}

type GetTotalExpensesForPeriodUseCase struct {
	PresentersRepository repositories.PresentersRepositoryInterface
	UserRepository       repositories.UserRepositoryInterface
}

func NewGetTotalExpensesForPeriodUseCase(
	PresentersRepository repositories.PresentersRepositoryInterface,
	UserRepository repositories.UserRepositoryInterface,
) *GetTotalExpensesForPeriodUseCase {
	return &GetTotalExpensesForPeriodUseCase{
		PresentersRepository: PresentersRepository,
		UserRepository:       UserRepository,
	}
}

func (c *GetTotalExpensesForPeriodUseCase) Execute(input GetTotalExpensesForPeriodInputDto) (GetTotalExpensesForPeriodOutputDto, []util.ProblemDetails) {
	user, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return GetTotalExpensesForPeriodOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
			},
		}
	} else if !user.Active {
		return GetTotalExpensesForPeriodOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Forbidden",
				Title:    "User is not active",
				Status:   403,
				Detail:   "User is not active",
				Instance: util.RFC403,
			},
		}
	}

	startDate, err := util.ParseDate(input.StartDate)
	if err != nil {
		return GetTotalExpensesForPeriodOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Bad Request",
				Title:    "Invalid start date",
				Status:   400,
				Detail:   "Start date is not in the correct format",
				Instance: util.RFC400,
			},
		}
	}

	endDate, err := util.ParseDate(input.EndDate)
	if err != nil {
		return GetTotalExpensesForPeriodOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Bad Request",
				Title:    "Invalid end date",
				Status:   400,
				Detail:   "End date is not in the correct format",
				Instance: util.RFC400,
			},
		}
	}

	if startDate.After(endDate) {
		return GetTotalExpensesForPeriodOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Bad Request",
				Title:    "Invalid date range",
				Status:   400,
				Detail:   "Start date must be before end date",
				Instance: util.RFC400,
			},
		}
	}

	total, err := c.PresentersRepository.GetTotalExpensesForPeriod(input.UserID, startDate, endDate)
	if err != nil {
		return GetTotalExpensesForPeriodOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Could not calculate total expenses",
				Status:   500,
				Detail:   err.Error(),
				Instance: util.RFC500,
			},
		}
	}

	return GetTotalExpensesForPeriodOutputDto{
		Total: total,
	}, nil
}
