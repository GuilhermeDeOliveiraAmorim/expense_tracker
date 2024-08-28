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
}
