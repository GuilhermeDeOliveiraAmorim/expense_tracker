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

type PresentersRepositoryInterface interface {
	GetTotalExpensesForPeriod(userID string, StartDate time.Time, EndDate time.Time) (float64, error)
	GetExpensesByCategoryPeriod(userID string, StartDate time.Time, EndDate time.Time) ([]CategoryExpense, error)
	GetMonthlyExpensesByCategoryYear(userID string, Year int) ([]MonthlyCategoryExpense, []int, error)
	GetMonthlyExpensesByTagYear(userID string, Year int) ([]MonthlyTagExpense, []int, error)
}
