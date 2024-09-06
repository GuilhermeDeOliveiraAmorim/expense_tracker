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
	CategoryID    string
	CategoryName  string
	CategoryColor string
	TotalAmount   float64
}, error) {
	var result []struct {
		CategoryID    string
		CategoryName  string
		CategoryColor string
		TotalAmount   float64
	}

	err := p.gorm.Table("expenses").
		Select("categories.id as category_id, categories.name as category_name, categories.color as category_color, SUM(expenses.amount) as total_amount").
		Joins("JOIN categories ON categories.id = expenses.category_id").
		Where("expenses.user_id = ? AND expenses.expanse_date BETWEEN ? AND ?", userID, periodStart, periodEnd).
		Group("categories.id, categories.name, categories.color").
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *PresentersRepository) ShowCategoryTreemapAmountPeriod(userID string, periodStart time.Time, periodEnd time.Time) ([]struct {
	CategoryName  string
	CategoryColor string
	TotalAmount   float64
}, error) {
	var result []struct {
		CategoryName  string
		CategoryColor string
		TotalAmount   float64
	}

	err := p.gorm.Table("expenses").
		Select("categories.name as category_name, categories.color as category_color, SUM(expenses.amount) as total_amount").
		Joins("JOIN categories ON categories.id = expenses.category_id").
		Where("expenses.user_id = ? AND expenses.expanse_date BETWEEN ? AND ?", userID, periodStart, periodEnd).
		Group("categories.name, categories.color").
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *PresentersRepository) ShowExpenseSimpleTablePeriod(userID string, periodStart time.Time, periodEnd time.Time, limit int, offset int) ([]struct {
	ExpenseID     string
	Amount        float64
	ExpenseDate   string
	Notes         string
	CategoryName  string
	CategoryColor string
}, error) {
	var result []struct {
		ExpenseID     string
		Amount        float64
		ExpenseDate   string
		Notes         string
		CategoryName  string
		CategoryColor string
	}

	err := p.gorm.Table("expenses").
		Select("expenses.id as expense_id, expenses.amount, expenses.expanse_date as expense_date, expenses.notes, categories.name as category_name, categories.color as category_color").
		Joins("JOIN categories ON categories.id = expenses.category_id").
		Where("expenses.user_id = ? AND expenses.expanse_date BETWEEN ? AND ?", userID, periodStart, periodEnd).
		Order("expenses.expanse_date DESC").
		Limit(limit).
		Offset(offset).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}
