package repositoriesgorm

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
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

func (p *PresentersRepository) GetMonthlyExpensesByCategoryYear(userID string, year int) ([]repositories.MonthlyCategoryExpense, []int, error) {
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
		Where("expenses.user_id = ? AND EXTRACT(YEAR FROM expenses.expanse_date) = ? AND expenses.active = ?", userID, year, true).
		Group("year, month, categories.name, categories.color").
		Order("MIN(expanse_date)").
		Scan(&results).Error

	if err != nil {
		tx.Rollback()
		return nil, []int{}, errors.New("failed to fetch monthly expenses by category: " + err.Error())
	}

	var years []int
	err = tx.Table("expenses").
		Select("DISTINCT EXTRACT(YEAR FROM expanse_date) AS year").
		Where("expenses.user_id = ? AND expenses.active = ?", userID, true).
		Order("year").
		Scan(&years).Error

	if err != nil {
		tx.Rollback()
		return nil, nil, errors.New("failed to fetch available years: " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return nil, []int{}, errors.New("failed to commit transaction")
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
				Total:      0,
			}
		}

		current := monthlyExpensesMap[key]

		current.Categories = append(current.Categories, repositories.CategoryExpense{
			CategoryName:  result.CategoryName,
			CategoryColor: result.Color,
			Total:         result.Total,
		})

		current.Total += result.Total

		monthlyExpensesMap[key] = current
	}

	var monthlyExpenses []repositories.MonthlyCategoryExpense
	for _, categories := range monthlyExpensesMap {
		monthlyExpenses = append(monthlyExpenses, categories)
	}

	sort.Slice(monthlyExpenses, func(i, j int) bool {
		return getMonthOrder(monthlyExpenses[i].Month) < getMonthOrder(monthlyExpenses[j].Month)
	})

	return monthlyExpenses, years, nil
}

func (p *PresentersRepository) GetMonthlyExpensesByTagYear(userID string, year int) ([]repositories.MonthlyTagExpense, []int, error) {
	tx := p.gorm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var results []struct {
		Year    int     `gorm:"column:year"`
		Month   string  `gorm:"column:month"`
		TagName string  `gorm:"column:tag_name"`
		Color   string  `gorm:"column:color"`
		Total   float64 `gorm:"column:total"`
	}

	err := tx.Table("expenses").
		Select("EXTRACT(YEAR FROM expanse_date) AS year, TO_CHAR(expanse_date, 'Month') AS month, tags.name AS tag_name, tags.color AS color, SUM(expenses.amount) AS total").
		Joins("INNER JOIN expense_tags ON expenses.id = expense_tags.expenses_id").
		Joins("INNER JOIN tags ON expense_tags.tags_id = tags.id").
		Where("expenses.user_id = ? AND EXTRACT(YEAR FROM expenses.expanse_date) = ? AND expenses.active = ?", userID, year, true).
		Group("year, month, tags.name, tags.color").
		Order("MIN(expanse_date)").
		Scan(&results).Error

	if err != nil {
		tx.Rollback()
		return nil, []int{}, errors.New("failed to fetch monthly expenses by tag: " + err.Error())
	}

	var years []int

	err = tx.Table("expenses").
		Select("DISTINCT EXTRACT(YEAR FROM expanse_date) AS year").
		Where("expenses.user_id = ? AND expenses.active = ?", userID, true).
		Order("year").
		Scan(&years).Error

	if err != nil {
		tx.Rollback()
		return nil, nil, errors.New("failed to fetch available years: " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return nil, []int{}, errors.New("failed to commit transaction")
	}

	monthlyExpensesMap := make(map[string]repositories.MonthlyTagExpense)

	for _, result := range results {
		month := strings.TrimSpace(result.Month)
		key := fmt.Sprintf("%d-%s", result.Year, month)

		if _, exists := monthlyExpensesMap[key]; !exists {
			monthlyExpensesMap[key] = repositories.MonthlyTagExpense{
				Month: month,
				Year:  result.Year,
				Tags:  []repositories.TagExpense{},
				Total: 0,
			}
		}

		current := monthlyExpensesMap[key]

		current.Tags = append(current.Tags, repositories.TagExpense{
			TagName:  result.TagName,
			TagColor: result.Color,
			Total:    result.Total,
		})

		current.Total += result.Total

		monthlyExpensesMap[key] = current
	}

	var monthlyExpenses []repositories.MonthlyTagExpense
	for _, tags := range monthlyExpensesMap {
		monthlyExpenses = append(monthlyExpenses, tags)
	}

	sort.Slice(monthlyExpenses, func(i, j int) bool {
		return getMonthOrder(monthlyExpenses[i].Month) < getMonthOrder(monthlyExpenses[j].Month)
	})

	return monthlyExpenses, years, nil
}

func getMonthOrder(month string) int {
	months := map[string]int{
		"January": 1, "February": 2, "March": 3, "April": 4, "May": 5, "June": 6,
		"July": 7, "August": 8, "September": 9, "October": 10, "November": 11, "December": 12,
	}
	return months[month]
}

func (p *PresentersRepository) GetTotalExpensesForCurrentMonth(userID string) (float64, string, error) {
	location, err := time.LoadLocation(util.TIMEZONE)
	if err != nil {
		return 0, "", errors.New("failed to load timezone: " + err.Error())
	}

	tx := p.gorm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var total float64
	var month string

	now := time.Now().In(location)
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, location)
	endOfMonth := now

	if err := tx.Table("expenses").
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ? AND expanse_date BETWEEN ? AND ? AND active = ?", userID, startOfMonth, endOfMonth, true).
		Scan(&total).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			total = 0
		} else {
			tx.Rollback()
			return 0, "", errors.New("failed to fetch total expenses: " + err.Error())
		}
	}

	month = time.Now().Format("January")

	if err := tx.Commit().Error; err != nil {
		return 0, "", errors.New("failed to commit transaction")
	}

	return total, month, nil
}

func (p *PresentersRepository) GetExpensesByMonthYear(userID string, month int, year int) (repositories.MonthExpenses, error) {
	tx := p.gorm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var monthExpenses repositories.MonthExpenses
	monthExpenses.Month = time.Month(month).String()
	monthExpenses.Year = year

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Nanosecond)

	var expenses []Expenses
	if err := tx.Where("user_id = ? AND expanse_date BETWEEN ? AND ? AND active = ?", userID, startDate, endDate, true).
		Find(&expenses).Error; err != nil {
		return repositories.MonthExpenses{}, errors.New("failed to fetch expenses: " + err.Error())
	}

	weeks := make(map[int]map[string]*repositories.DayExpense)

	for _, expense := range expenses {

		_, weekNumber := expense.ExpanseDate.ISOWeek()

		weekNumber = int(weekNumber)

		dayKey := expense.ExpanseDate.Format("02")

		if weeks[weekNumber] == nil {
			weeks[weekNumber] = make(map[string]*repositories.DayExpense)
		}

		if _, exists := weeks[weekNumber][dayKey]; !exists {
			weeks[weekNumber][dayKey] = &repositories.DayExpense{
				Day:     dayKey,
				DayName: expense.ExpanseDate.Weekday().String(),
				Total:   0,
				Tags:    []repositories.ExpenseTag{},
			}
		}

		var tags []Tags
		if err := tx.Table("tags").
			Joins("JOIN expense_tags ON tags.id = expense_tags.tags_id").
			Where("expense_tags.expenses_id = ?", expense.ID).
			Select("tags.name, tags.color").
			Find(&tags).Error; err != nil {
			return repositories.MonthExpenses{}, errors.New("failed to fetch tags for expense: " + err.Error())
		}

		weeks[weekNumber][dayKey].Total += expense.Amount

		for _, tag := range tags {
			tagFound := false

			for i, dayTag := range weeks[weekNumber][dayKey].Tags {
				if dayTag.Name == tag.Name {
					weeks[weekNumber][dayKey].Tags[i].Total += expense.Amount
					tagFound = true
					break
				}
			}

			if !tagFound {
				weeks[weekNumber][dayKey].Tags = append(weeks[weekNumber][dayKey].Tags, repositories.ExpenseTag{
					Name:  tag.Name,
					Color: tag.Color,
					Total: expense.Amount,
				})
			}
		}
	}

	for weekNumber, days := range weeks {
		weekExpenses := repositories.WeekExpenses{
			Week: weekNumber,
			Days: []repositories.DayExpense{},
		}

		for _, day := range days {
			weekExpenses.Days = append(weekExpenses.Days, *day)
		}

		sort.Slice(weekExpenses.Days, func(i, j int) bool {
			dayI, _ := strconv.Atoi(weekExpenses.Days[i].Day)
			dayJ, _ := strconv.Atoi(weekExpenses.Days[j].Day)
			return dayI < dayJ
		})

		monthExpenses.Weeks = append(monthExpenses.Weeks, weekExpenses)
	}

	sort.Slice(monthExpenses.Weeks, func(i, j int) bool {
		return monthExpenses.Weeks[i].Week < monthExpenses.Weeks[j].Week
	})

	return monthExpenses, nil
}

func (p *PresentersRepository) GetTotalExpensesForCurrentWeek(userID string) (float64, string, error) {
	tx := p.gorm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var totalExpenses float64

	now := time.Now()

	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	endOfWeek := now

	var expenses []Expenses
	if err := p.gorm.
		Where("user_id = ? AND expanse_date BETWEEN ? AND ?", userID, startOfMonth, endOfWeek).
		Find(&expenses).Error; err != nil {
		return 0, "", errors.New("failed to fetch expenses: " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return 0, "", errors.New("failed to commit transaction")
	}

	for _, expense := range expenses {
		totalExpenses += expense.Amount
	}

	weekInterval := fmt.Sprintf("%s - %s", startOfMonth.Format("02/01/2006"), endOfWeek.Format("02/01/2006"))

	return totalExpenses, weekInterval, nil
}
