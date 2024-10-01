package repositories

import (
	"time"
)

type CategoryExpense struct {
	CategoryName  string  `json:"category_name"`
	CategoryColor string  `json:"category_color"`
	Total         float64 `json:"total"`
}

type MonthlyCategoryExpense struct {
	Month      string            `json:"month"`
	Year       int               `json:"year"`
	Categories []CategoryExpense `json:"categories"`
}

type PresentersRepositoryInterface interface {
	GetTotalExpensesForPeriod(userID string, StartDate time.Time, EndDate time.Time) (float64, error)
	GetExpensesByCategoryPeriod(userID string, StartDate time.Time, EndDate time.Time) ([]CategoryExpense, error)
	GetMonthlyExpensesByCategoryPeriod(userID string, Year int) ([]MonthlyCategoryExpense, []int, error)
}
