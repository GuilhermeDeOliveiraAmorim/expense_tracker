package repositories

import (
	"time"
)

type PresentersRepositoryInterface interface {
	ShowTotalExpensesCategoryPeriod(userID string, periodStart time.Time, periodEnd time.Time) ([]struct {
		CategoryID   string
		CategoryName string
		TotalAmount  float64
	}, error)
}
