package repositoriesgorm

import (
	"time"

	"gorm.io/gorm"
)

type PresentersRepository struct {
	gorm *gorm.DB
}

func NewPresentersRepository(gorm *gorm.DB) *PresentersRepository {
	return &PresentersRepository{
		gorm: gorm,
	}
}

func (p *PresentersRepository) ShowTotalExpensesCategoryPeriod(userID string, periodStart time.Time, periodEnd time.Time) ([]struct {
	CategoryID   string
	CategoryName string
	TotalAmount  float64
}, error) {
	var result []struct {
		CategoryID   string
		CategoryName string
		TotalAmount  float64
	}

	err := p.gorm.Table("expenses").
		Select("categories.name as category_name, SUM(expenses.amount) as total_amount").
		Joins("JOIN categories ON categories.id = expenses.category_id").
		Where("expenses.user_id = ? AND expenses.expanse_date BETWEEN ? AND ?", userID, periodStart, periodEnd).
		Group("categories.name").
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}
