package repositories

import (
	"time"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
)

type CategoryExpense struct {
	CategoryName  string  `json:"category_name"`
	CategoryColor string  `json:"category_color"`
	Total         float64 `json:"total"`
}

type TagExpense struct {
	TagName  string  `json:"tag_name"`
	TagColor string  `json:"tag_color"`
	Total    float64 `json:"total"`
}

type MonthlyCategoryExpense struct {
	Month      string            `json:"month"`
	Year       int               `json:"year"`
	Categories []CategoryExpense `json:"categories"`
	Total      float64           `json:"total"`
}

type MonthlyTagExpense struct {
	Month string       `json:"month"`
	Year  int          `json:"year"`
	Tags  []TagExpense `json:"tags"`
	Total float64      `json:"total"`
}

type MonthExpenses struct {
	Month          string         `json:"month"`
	Year           int            `json:"year"`
	TotalExpenses  float64        `json:"total_expenses"`
	Weeks          []WeekExpenses `json:"weeks"`
	AvailableYears []int          `json:"available_years"`
}

type WeekExpenses struct {
	Week int          `json:"week"`
	Days []DayExpense `json:"days"`
}

type DayExpense struct {
	Day     string       `json:"day"`
	DayName string       `json:"day_name"`
	Total   float64      `json:"total"`
	Tags    []ExpenseTag `json:"tags"`
}

type ExpenseTag struct {
	Name  string  `json:"name"`
	Color string  `json:"color"`
	Total float64 `json:"total"`
}

type MonthCurrentYear struct {
	Month string  `json:"month"`
	Total float64 `json:"total"`
}

type ExpensesMonthCurrentYear struct {
	Year           int                `json:"year"`
	Total          float64            `json:"total"`
	Months         []MonthCurrentYear `json:"months"`
	AvailableYears []int              `json:"available_years"`
}

type CategoryTagTotal struct {
	Name      string  `json:"name"`
	TagAmount float64 `json:"tag_amount"`
}

type CategoryWithTags struct {
	Name           string             `json:"name"`
	CategoryAmount float64            `json:"category_amount"`
	Tags           []CategoryTagTotal `json:"tags"`
}

type MonthOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type CategoryTagsTotals struct {
	Month           string             `json:"month"`
	Year            int                `json:"year"`
	ExpensesAmount  float64            `json:"expenses_amount"`
	Categories      []CategoryWithTags `json:"categories"`
	AvailableYears  []int              `json:"available_years"`
	AvailableMonths []MonthOption      `json:"available_months"`
}

type PresentersRepositoryInterface interface {
	GetTotalExpensesForPeriod(userID string, StartDate time.Time, EndDate time.Time) (float64, error)
	GetExpensesByCategoryPeriod(userID string, StartDate time.Time, EndDate time.Time) ([]CategoryExpense, error)
	GetMonthlyExpensesByCategoryYear(userID string, Year int) ([]MonthlyCategoryExpense, []int, error)
	GetMonthlyExpensesByTagYear(userID string, Year int) ([]MonthlyTagExpense, []int, error)
	GetTotalExpensesForCurrentMonth(userID string) (float64, string, error)
	GetExpensesByMonthYear(userID string, month int, year int) (MonthExpenses, error)
	GetTotalExpensesForCurrentWeek(userID string) (float64, string, error)
	GetTotalExpensesMonthCurrentYear(userID string, year int) (ExpensesMonthCurrentYear, error)
	GetCategoryTagsTotalsByMonthYear(userID string, month int, year int) (CategoryTagsTotals, error)
	GetAvailableMonthsYears(userID string) ([]int, []MonthOption, error)
	GetDayToDayExpensesPeriod(userID string, StartDate time.Time, EndDate time.Time) ([]entities.Expense, error)
}
