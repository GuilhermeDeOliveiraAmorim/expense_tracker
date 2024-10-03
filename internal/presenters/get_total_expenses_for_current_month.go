package presenters

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type GetTotalExpensesForCurrentMonthInputDto struct {
	UserID string `json:"user_id"`
}

type GetTotalExpensesForCurrentMonthOutputDto struct {
	TotalExpenses float64 `json:"total_expenses"`
	CurrentMonth  string  `json:"current_month"`
}

type GetTotalExpensesForCurrentMonthUseCase struct {
	PresentersRepository repositories.PresentersRepositoryInterface
	UserRepository       repositories.UserRepositoryInterface
}

func NewGetTotalExpensesForCurrentMonthUseCase(
	PresentersRepository repositories.PresentersRepositoryInterface,
	UserRepository repositories.UserRepositoryInterface,
) *GetTotalExpensesForCurrentMonthUseCase {
	return &GetTotalExpensesForCurrentMonthUseCase{
		PresentersRepository: PresentersRepository,
		UserRepository:       UserRepository,
	}
}

func (c *GetTotalExpensesForCurrentMonthUseCase) Execute(input GetTotalExpensesForCurrentMonthInputDto) (GetTotalExpensesForCurrentMonthOutputDto, []util.ProblemDetails) {
	user, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return GetTotalExpensesForCurrentMonthOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
			},
		}
	} else if !user.Active {
		return GetTotalExpensesForCurrentMonthOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Forbidden",
				Title:    "User is not active",
				Status:   403,
				Detail:   "User is not active",
				Instance: util.RFC403,
			},
		}
	}

	total, month, err := c.PresentersRepository.GetTotalExpensesForCurrentMonth(input.UserID)
	if err != nil {
		return GetTotalExpensesForCurrentMonthOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Could not calculate total expenses",
				Status:   500,
				Detail:   err.Error(),
				Instance: util.RFC500,
			},
		}
	}

	return GetTotalExpensesForCurrentMonthOutputDto{
		TotalExpenses: total,
		CurrentMonth:  month,
	}, nil
}
