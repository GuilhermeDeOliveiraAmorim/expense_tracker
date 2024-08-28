package presenters

import (
	"time"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type NameAmount struct {
	CategoryName string  `json:"x"`
	TotalAmount  float64 `json:"y"`
}

type CategoryTreemapAmountPeriod struct {
	NameAmount []NameAmount `json:"data"`
	Colors     []string     `json:"colors"`
}

type ShowCategoryTreemapAmountPeriodInputDto struct {
	UserID      string `json:"user_id"`
	PeriodStart string `json:"period_start"`
	PeriodEnd   string `json:"period_end"`
}

type ShowCategoryTreemapAmountPeriodOutputDto struct {
	CategoryTreemapAmountPeriod CategoryTreemapAmountPeriod `json:"category_treemap_amount_period"`
}

type ShowCategoryTreemapAmountPeriodPresenters struct {
	PresentersRepository repositories.PresentersRepositoryInterface
	UserRepository       repositories.UserRepositoryInterface
}

func NewCategoryTreemapAmountPeriodPresenters(
	PresentersRepository repositories.PresentersRepositoryInterface,
	UserRepository repositories.UserRepositoryInterface,
) *ShowCategoryTreemapAmountPeriodPresenters {
	return &ShowCategoryTreemapAmountPeriodPresenters{
		PresentersRepository: PresentersRepository,
		UserRepository:       UserRepository,
	}
}

func (p *ShowCategoryTreemapAmountPeriodPresenters) Execute(input ShowCategoryTreemapAmountPeriodInputDto) (ShowCategoryTreemapAmountPeriodOutputDto, []util.ProblemDetails) {
	user, err := p.UserRepository.GetUser(input.UserID)
	if err != nil {
		return ShowCategoryTreemapAmountPeriodOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
			},
		}
	} else if !user.Active {
		return ShowCategoryTreemapAmountPeriodOutputDto{}, []util.ProblemDetails{
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
		return ShowCategoryTreemapAmountPeriodOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Validation Error",
				Title:    "Bad Request",
				Status:   400,
				Detail:   "Invalid timezone",
				Instance: util.RFC400,
			},
		}
	}

	periodStart, err := time.ParseInLocation(util.DATEFORMAT, input.PeriodStart, location)
	if err != nil {
		return ShowCategoryTreemapAmountPeriodOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Validation Error",
				Title:    "Bad Request",
				Status:   400,
				Detail:   "Invalid start date format",
				Instance: util.RFC400,
			},
		}
	}

	periodEnd, err := time.ParseInLocation(util.DATEFORMAT, input.PeriodEnd, location)
	if err != nil {
		return ShowCategoryTreemapAmountPeriodOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Validation Error",
				Title:    "Bad Request",
				Status:   400,
				Detail:   "Invalid end date format",
				Instance: util.RFC400,
			},
		}
	}

	totalExpensesByCategory, err := p.PresentersRepository.ShowCategoryTreemapAmountPeriod(user.ID, periodStart, periodEnd)
	if err != nil {
		return ShowCategoryTreemapAmountPeriodOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "An error occurred while retrieving total expenses by category",
				Status:   500,
				Detail:   err.Error(),
				Instance: util.RFC400,
			},
		}
	}

	if len(totalExpensesByCategory) == 0 {
		return ShowCategoryTreemapAmountPeriodOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "No expenses found within the specified period or user",
				Status:   404,
				Detail:   "No expenses found within the specified period or user",
				Instance: util.RFC404,
			},
		}
	}

	var nameAmounts []NameAmount
	var colors []string

	for _, expense := range totalExpensesByCategory {
		nameAmounts = append(nameAmounts, NameAmount{
			CategoryName: expense.CategoryName,
			TotalAmount:  expense.TotalAmount,
		})
		colors = append(colors, expense.CategoryColor)
	}

	return ShowCategoryTreemapAmountPeriodOutputDto{
		CategoryTreemapAmountPeriod: CategoryTreemapAmountPeriod{
			NameAmount: nameAmounts,
			Colors:     colors,
		},
	}, nil
}
