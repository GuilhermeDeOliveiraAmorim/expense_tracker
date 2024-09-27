package repositoriesgorm

import (
	"errors"
	"time"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
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

func (p *PresentersRepository) GetTotalExpensesForPeriod(userID string, startDate time.Time, endDate time.Time) (float64, error) {
	tx := p.gorm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var total float64

	if err := tx.Table("expenses").
		Select("SUM(amount)").
		Where("user_id = ? AND expanse_date BETWEEN ? AND ? AND active = ?", userID, startDate, endDate, true).
		Scan(&total).Error; err != nil {
		tx.Rollback()
		return 0, errors.New("failed to calculate total expenses: " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return 0, errors.New("failed to commit transaction")
	}

	return total, nil
}

func (p *PresentersRepository) GetExpensesByCategoryPeriod(userID string, startDate time.Time, endDate time.Time) ([]repositories.CategoryExpense, error) {
	tx := p.gorm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var expensesByCategory []repositories.CategoryExpense

	if err := tx.Table("expenses").
		Select("categories.name as category_name, categories.color as category_color, SUM(expenses.amount) as total").
		Joins("JOIN categories ON expenses.category_id = categories.id").
		Where("expenses.user_id = ? AND expenses.expanse_date BETWEEN ? AND ? AND expenses.active = ?", userID, startDate, endDate, true).
		Group("categories.name, categories.color").
		Scan(&expensesByCategory).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to fetch expenses by category: " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("failed to commit transaction")
	}

	return expensesByCategory, nil
}
