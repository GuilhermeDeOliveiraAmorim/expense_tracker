package repositories

import (
	"time"
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
	Month         string         `json:"month"`
	Year          int            `json:"year"`
	TotalExpenses float64        `json:"total_expenses"`
	Weeks         []WeekExpenses `json:"weeks"`
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

type PresentersRepositoryInterface interface {
	GetTotalExpensesForPeriod(userID string, StartDate time.Time, EndDate time.Time) (float64, error)
	GetExpensesByCategoryPeriod(userID string, StartDate time.Time, EndDate time.Time) ([]CategoryExpense, error)
	GetMonthlyExpensesByCategoryYear(userID string, Year int) ([]MonthlyCategoryExpense, []int, error)
	GetMonthlyExpensesByTagYear(userID string, Year int) ([]MonthlyTagExpense, []int, error)
	GetTotalExpensesForCurrentMonth(userID string) (float64, string, error)
	GetExpensesByMonthYear(userID string, month int, year int) (MonthExpenses, error)
	GetTotalExpensesForCurrentWeek(userID string) (float64, string, error)
	GetTotalExpensesMonthCurrentYear(userID string, year int) (ExpensesMonthCurrentYear, error)
}
