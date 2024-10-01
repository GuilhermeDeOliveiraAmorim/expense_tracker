package repositoriesgorm

import (
	"errors"
	"fmt"
	"sort"
	"strings"
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
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ? AND expanse_date BETWEEN ? AND ? AND active = ?", userID, startDate, endDate, true).
		Scan(&total).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			total = 0
		} else {
			tx.Rollback()
			return 0, errors.New("failed to fetch total expenses: " + err.Error())
		}
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
		Group("categories.name, categories.color").Order("total").
		Scan(&expensesByCategory).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to fetch expenses by category: " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("failed to commit transaction")
	}

	return expensesByCategory, nil
}

func (p *PresentersRepository) GetMonthlyExpensesByCategoryPeriod(userID string, startDate time.Time, endDate time.Time) ([]repositories.MonthlyCategoryExpense, error) {
	tx := p.gorm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var results []struct {
		Year         int     `gorm:"column:year"`
		Month        string  `gorm:"column:month"`
		CategoryName string  `gorm:"column:category_name"`
		Color        string  `gorm:"column:color"`
		Total        float64 `gorm:"column:total"`
	}

	err := tx.Table("expenses").
		Select("EXTRACT(YEAR FROM expanse_date) AS year, TO_CHAR(expanse_date, 'Month') AS month, categories.name AS category_name, categories.color AS color, SUM(expenses.amount) AS total").
		Joins("INNER JOIN categories ON expenses.category_id = categories.id").
		Where("expenses.user_id = ? AND expenses.expanse_date BETWEEN ? AND ? AND expenses.active = ?", userID, startDate, endDate, true).
		Group("year, month, categories.name, categories.color").
		Order("MIN(expanse_date)").
		Scan(&results).Error

	if err != nil {
		tx.Rollback()
		return nil, errors.New("failed to fetch monthly expenses by category: " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("failed to commit transaction")
	}

	monthlyExpensesMap := make(map[string]repositories.MonthlyCategoryExpense)

	for _, result := range results {
		month := strings.TrimSpace(result.Month)
		key := fmt.Sprintf("%d-%s", result.Year, month)

		if _, exists := monthlyExpensesMap[key]; !exists {
			monthlyExpensesMap[key] = repositories.MonthlyCategoryExpense{
				Month:      month,
				Year:       result.Year,
				Categories: []repositories.CategoryExpense{},
			}
		}

		current := monthlyExpensesMap[key]

		current.Categories = append(current.Categories, repositories.CategoryExpense{
			CategoryName:  result.CategoryName,
			CategoryColor: result.Color,
			Total:         result.Total,
		})

		monthlyExpensesMap[key] = current
	}

	var monthlyExpenses []repositories.MonthlyCategoryExpense
	for _, categories := range monthlyExpensesMap {
		monthlyExpenses = append(monthlyExpenses, categories)
	}

	sort.Slice(monthlyExpenses, func(i, j int) bool {
		return getMonthOrder(monthlyExpenses[i].Month) < getMonthOrder(monthlyExpenses[j].Month)
	})

	return monthlyExpenses, nil
}

func getMonthOrder(month string) int {
	months := map[string]int{
		"January": 1, "February": 2, "March": 3, "April": 4, "May": 5, "June": 6,
		"July": 7, "August": 8, "September": 9, "October": 10, "November": 11, "December": 12,
	}
	return months[month]
}
