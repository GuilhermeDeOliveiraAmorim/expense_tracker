package repositories

import (
	"time"
)

type PresentersRepositoryInterface interface {
	ShowTotalExpensesCategoryPeriod(userID string, periodStart time.Time, periodEnd time.Time) ([]struct {
		CategoryID    string
		CategoryName  string
		CategoryColor string
		TotalAmount   float64
	}, error)
	ShowCategoryTreemapAmountPeriod(userID string, periodStart time.Time, periodEnd time.Time) ([]struct {
		CategoryName  string
		CategoryColor string
		TotalAmount   float64
	}, error)
	ShowExpenseSimpleTablePeriod(userID string, periodStart time.Time, periodEnd time.Time, limit int, offset int) ([]struct {
		ExpenseID     string
		Amount        float64
		ExpenseDate   string
		Notes         string
		CategoryName  string
		CategoryColor string
	}, error)
}
