package repositories

import (
	"time"
)

type CategoryExpense struct {
	CategoryName  string  `json:"category_name"`
	CategoryColor string  `json:"category_color"`
	Total         float64 `json:"total"`
}

type PresentersRepositoryInterface interface {
	GetTotalExpensesForPeriod(userID string, StartDate time.Time, EndDate time.Time) (float64, error)
	GetExpensesByCategoryPeriod(userID string, StartDate time.Time, EndDate time.Time) ([]CategoryExpense, error)
}
