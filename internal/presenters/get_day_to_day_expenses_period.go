package presenters

import (
	"strconv"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type DayToDayExpense struct {
	Day     string  `json:"day"`
	DayName string  `json:"day_name"`
	Month   string  `json:"month"`
	Year    string  `json:"year"`
	Amount  float64 `json:"amount"`
}

type GetDayToDayExpensesPeriodInputDto struct {
	UserID    string `json:"user_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type GetDayToDayExpensesPeriodOutputDto struct {
	Expenses []DayToDayExpense `json:"expenses"`
}

type GetDayToDayExpensesPeriodUseCase struct {
	PresentersRepository repositories.PresentersRepositoryInterface
	UserRepository       repositories.UserRepositoryInterface
}

func NewGetDayToDayExpensesPeriodUseCase(
	PresentersRepository repositories.PresentersRepositoryInterface,
	UserRepository repositories.UserRepositoryInterface,
) *GetDayToDayExpensesPeriodUseCase {
	return &GetDayToDayExpensesPeriodUseCase{
		PresentersRepository: PresentersRepository,
		UserRepository:       UserRepository,
	}
}

func (c *GetDayToDayExpensesPeriodUseCase) Execute(input GetDayToDayExpensesPeriodInputDto) (GetDayToDayExpensesPeriodOutputDto, []util.ProblemDetails) {
	user, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return GetDayToDayExpensesPeriodOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
			},
		}
	} else if !user.Active {
		return GetDayToDayExpensesPeriodOutputDto{}, []util.ProblemDetails{
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
		return GetDayToDayExpensesPeriodOutputDto{}, []util.ProblemDetails{
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
		return GetDayToDayExpensesPeriodOutputDto{}, []util.ProblemDetails{
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
		return GetDayToDayExpensesPeriodOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Bad Request",
				Title:    "Invalid date range",
				Status:   400,
				Detail:   "Start date must be before end date",
				Instance: util.RFC400,
			},
		}
	}

	expenses, err := c.PresentersRepository.GetDayToDayExpensesPeriod(input.UserID, startDate, endDate)
	if err != nil {
		return GetDayToDayExpensesPeriodOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Could not calculate expenses day by day",
				Status:   500,
				Detail:   err.Error(),
				Instance: util.RFC500,
			},
		}
	}

	var dayToDayExpenses []DayToDayExpense

	for _, expense := range expenses {
		dayToDayExpenses = append(dayToDayExpenses, DayToDayExpense{
			Day:     expense.ExpenseDate.Format("02"),
			DayName: expense.ExpenseDate.Weekday().String(),
			Month:   expense.ExpenseDate.Month().String(),
			Year:    strconv.Itoa(expense.ExpenseDate.Year()),
			Amount:  expense.Amount,
		})
	}

	return GetDayToDayExpensesPeriodOutputDto{
		Expenses: dayToDayExpenses,
	}, nil
}
