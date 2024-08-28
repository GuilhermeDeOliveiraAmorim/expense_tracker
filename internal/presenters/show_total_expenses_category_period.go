package presenters

import (
	"time"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type TotalExpensesCategoryPeriod struct {
	CategoryID    string  `json:"category_id"`
	CategoryName  string  `json:"category_name"`
	CategoryColor string  `json:"category_color"`
	TotalAmount   float64 `json:"total_amount"`
}

type ShowTotalExpensesCategoryPeriodInputDto struct {
	UserID      string `json:"user_id"`
	PeriodStart string `json:"period_start"`
	PeriodEnd   string `json:"period_end"`
}

type ShowTotalExpensesCategoryPeriodOutputDto struct {
	TotalExpensesByCategory []TotalExpensesCategoryPeriod `json:"total_expenses_by_category"`
}

type ShowTotalExpensesCategoryPeriodPresenters struct {
	PresentersRepository repositories.PresentersRepositoryInterface
	UserRepository       repositories.UserRepositoryInterface
}

func NewShowTotalExpensesCategoryPeriodPresenters(
	PresentersRepository repositories.PresentersRepositoryInterface,
	UserRepository repositories.UserRepositoryInterface,
) *ShowTotalExpensesCategoryPeriodPresenters {
	return &ShowTotalExpensesCategoryPeriodPresenters{
		PresentersRepository: PresentersRepository,
		UserRepository:       UserRepository,
	}
}

func (p *ShowTotalExpensesCategoryPeriodPresenters) Execute(input ShowTotalExpensesCategoryPeriodInputDto) (ShowTotalExpensesCategoryPeriodOutputDto, []util.ProblemDetails) {
	user, err := p.UserRepository.GetUser(input.UserID)
	if err != nil {
		return ShowTotalExpensesCategoryPeriodOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
			},
		}
	} else if !user.Active {
		return ShowTotalExpensesCategoryPeriodOutputDto{}, []util.ProblemDetails{
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
		return ShowTotalExpensesCategoryPeriodOutputDto{}, []util.ProblemDetails{
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
		return ShowTotalExpensesCategoryPeriodOutputDto{}, []util.ProblemDetails{
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
		return ShowTotalExpensesCategoryPeriodOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Validation Error",
				Title:    "Bad Request",
				Status:   400,
				Detail:   "Invalid end date format",
				Instance: util.RFC400,
			},
		}
	}

	totalExpensesByCategory, err := p.PresentersRepository.ShowTotalExpensesCategoryPeriod(user.ID, periodStart, periodEnd)
	if err != nil {
		return ShowTotalExpensesCategoryPeriodOutputDto{}, []util.ProblemDetails{
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
		return ShowTotalExpensesCategoryPeriodOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "No expenses found within the specified period or user",
				Status:   404,
				Detail:   "No expenses found within the specified period or user",
				Instance: util.RFC404,
			},
		}
	}

	var expensesByCategory []TotalExpensesCategoryPeriod

	for _, expense := range totalExpensesByCategory {
		expensesByCategory = append(expensesByCategory, TotalExpensesCategoryPeriod{
			CategoryID:    expense.CategoryID,
			CategoryName:  expense.CategoryName,
			CategoryColor: expense.CategoryColor,
			TotalAmount:   expense.TotalAmount,
		})
	}

	return ShowTotalExpensesCategoryPeriodOutputDto{
		TotalExpensesByCategory: expensesByCategory,
	}, nil
}
