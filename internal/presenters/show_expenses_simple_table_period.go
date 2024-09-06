package presenters

import (
	"time"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type Row struct {
	ExpenseID     string  `json:"expense_id"`
	Amount        float64 `json:"amount"`
	ExpenseDate   string  `json:"expense_date"`
	Notes         string  `json:"notes"`
	CategoryName  string  `json:"category_name"`
	CategoryColor string  `json:"category_color"`
}

type ExpenseSimpleTable struct {
	NumberOfRows int   `json:"number_of_rows"`
	Rows         []Row `json:"rows"`
}

type ShowExpenseSimpleTablePeriodInputDto struct {
	UserID      string `json:"user_id"`
	PeriodStart string `json:"period_start"`
	PeriodEnd   string `json:"period_end"`
	Limit       int    `json:"limit,string"`
	Offset      int    `json:"offset,string"`
}

type ShowExpenseSimpleTablePeriodOutputDto struct {
	ExpenseSimpleTable ExpenseSimpleTable `json:"expense_simple_table"`
}

type ShowExpenseSimpleTablePeriodPresenters struct {
	PresentersRepository repositories.PresentersRepositoryInterface
	UserRepository       repositories.UserRepositoryInterface
}

func NewShowExpenseSimpleTablePeriodPresenters(
	PresentersRepository repositories.PresentersRepositoryInterface,
	UserRepository repositories.UserRepositoryInterface,
) *ShowExpenseSimpleTablePeriodPresenters {
	return &ShowExpenseSimpleTablePeriodPresenters{
		PresentersRepository: PresentersRepository,
		UserRepository:       UserRepository,
	}
}

func (p *ShowExpenseSimpleTablePeriodPresenters) Execute(input ShowExpenseSimpleTablePeriodInputDto) (ShowExpenseSimpleTablePeriodOutputDto, []util.ProblemDetails) {
	user, err := p.UserRepository.GetUser(input.UserID)
	if err != nil {
		return ShowExpenseSimpleTablePeriodOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
			},
		}
	} else if !user.Active {
		return ShowExpenseSimpleTablePeriodOutputDto{}, []util.ProblemDetails{
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
		return ShowExpenseSimpleTablePeriodOutputDto{}, []util.ProblemDetails{
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
		return ShowExpenseSimpleTablePeriodOutputDto{}, []util.ProblemDetails{
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
		return ShowExpenseSimpleTablePeriodOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Validation Error",
				Title:    "Bad Request",
				Status:   400,
				Detail:   "Invalid end date format",
				Instance: util.RFC400,
			},
		}
	}

	dataTable, err := p.PresentersRepository.ShowExpenseSimpleTablePeriod(user.ID, periodStart, periodEnd, input.Limit, input.Offset)
	if err != nil {
		return ShowExpenseSimpleTablePeriodOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "An error occurred while retrieving total expenses by category",
				Status:   500,
				Detail:   err.Error(),
				Instance: util.RFC400,
			},
		}
	}

	if len(dataTable) == 0 {
		return ShowExpenseSimpleTablePeriodOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "No expenses found within the specified period or user",
				Status:   404,
				Detail:   "No expenses found within the specified period or user",
				Instance: util.RFC404,
			},
		}
	}

	var rows []Row

	for _, row := range dataTable {
		rows = append(rows, Row{
			ExpenseID:     row.ExpenseID,
			Amount:        row.Amount,
			ExpenseDate:   row.ExpenseDate,
			Notes:         row.Notes,
			CategoryName:  row.CategoryName,
			CategoryColor: row.CategoryColor,
		})
	}

	return ShowExpenseSimpleTablePeriodOutputDto{
		ExpenseSimpleTable: ExpenseSimpleTable{
			NumberOfRows: len(rows),
			Rows:         rows,
		},
	}, nil
}
