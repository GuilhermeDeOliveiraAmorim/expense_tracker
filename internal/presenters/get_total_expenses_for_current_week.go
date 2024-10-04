package presenters

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type GetTotalExpensesForCurrentWeekInputDto struct {
	UserID string `json:"user_id"`
}

type GetTotalExpensesForCurrentWeekOutputDto struct {
	TotalExpenses float64 `json:"total_expenses"`
	CurrentWeek   string  `json:"current_month"`
}

type GetTotalExpensesForCurrentWeekUseCase struct {
	PresentersRepository repositories.PresentersRepositoryInterface
	UserRepository       repositories.UserRepositoryInterface
}

func NewGetTotalExpensesForCurrentWeekUseCase(
	PresentersRepository repositories.PresentersRepositoryInterface,
	UserRepository repositories.UserRepositoryInterface,
) *GetTotalExpensesForCurrentWeekUseCase {
	return &GetTotalExpensesForCurrentWeekUseCase{
		PresentersRepository: PresentersRepository,
		UserRepository:       UserRepository,
	}
}

func (c *GetTotalExpensesForCurrentWeekUseCase) Execute(input GetTotalExpensesForCurrentWeekInputDto) (GetTotalExpensesForCurrentWeekOutputDto, []util.ProblemDetails) {
	user, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return GetTotalExpensesForCurrentWeekOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
			},
		}
	} else if !user.Active {
		return GetTotalExpensesForCurrentWeekOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Forbidden",
				Title:    "User is not active",
				Status:   403,
				Detail:   "User is not active",
				Instance: util.RFC403,
			},
		}
	}

	total, month, err := c.PresentersRepository.GetTotalExpensesForCurrentWeek(input.UserID)
	if err != nil {
		return GetTotalExpensesForCurrentWeekOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Could not calculate total expenses",
				Status:   500,
				Detail:   err.Error(),
				Instance: util.RFC500,
			},
		}
	}

	return GetTotalExpensesForCurrentWeekOutputDto{
		TotalExpenses: total,
		CurrentWeek:   month,
	}, nil
}
