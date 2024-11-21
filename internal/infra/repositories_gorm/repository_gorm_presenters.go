package repositoriesgorm

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
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
	var total float64

	if err := p.gorm.Table("expenses").
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ? AND expanse_date BETWEEN ? AND ? AND active = ?", userID, startDate, endDate, true).
		Scan(&total).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			total = 0
		} else {
			return 0, errors.New("failed to fetch total expenses: " + err.Error())
		}
	}

	return total, nil
}

func (p *PresentersRepository) GetExpensesByCategoryPeriod(userID string, startDate time.Time, endDate time.Time) ([]repositories.CategoryExpense, error) {
	var expensesByCategory []repositories.CategoryExpense

	if err := p.gorm.Table("expenses").
		Select("categories.name as category_name, categories.color as category_color, SUM(expenses.amount) as total").
		Joins("JOIN categories ON expenses.category_id = categories.id").
		Where("expenses.user_id = ? AND expenses.expanse_date BETWEEN ? AND ? AND expenses.active = ?", userID, startDate, endDate, true).
		Group("categories.name, categories.color").Order("total DESC").
		Scan(&expensesByCategory).Error; err != nil {
		return nil, errors.New("failed to fetch expenses by category: " + err.Error())
	}

	return expensesByCategory, nil
}

func (p *PresentersRepository) GetMonthlyExpensesByCategoryYear(userID string, year int) ([]repositories.MonthlyCategoryExpense, []int, error) {
	var results []struct {
		Year         int     `gorm:"column:year"`
		Month        string  `gorm:"column:month"`
		CategoryName string  `gorm:"column:category_name"`
		Color        string  `gorm:"column:color"`
		Total        float64 `gorm:"column:total"`
	}

	err := p.gorm.Table("expenses").
		Select("EXTRACT(YEAR FROM expanse_date) AS year, TO_CHAR(expanse_date, 'Month') AS month, categories.name AS category_name, categories.color AS color, SUM(expenses.amount) AS total").
		Joins("INNER JOIN categories ON expenses.category_id = categories.id").
		Where("expenses.user_id = ? AND EXTRACT(YEAR FROM expenses.expanse_date) = ? AND expenses.active = ?", userID, year, true).
		Group("year, month, categories.name, categories.color").
		Order("MIN(expanse_date)").
		Scan(&results).Error

	if err != nil {
		return nil, []int{}, errors.New("failed to fetch monthly expenses by category: " + err.Error())
	}

	var years []int
	err = p.gorm.Table("expenses").
		Select("DISTINCT EXTRACT(YEAR FROM expanse_date) AS year").
		Where("expenses.user_id = ? AND expenses.active = ?", userID, true).
		Order("year").
		Scan(&years).Error

	if err != nil {
		p.gorm.Rollback()
		return nil, nil, errors.New("failed to fetch available years: " + err.Error())
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
	var results []struct {
		Year    int     `gorm:"column:year"`
		Month   string  `gorm:"column:month"`
		TagName string  `gorm:"column:tag_name"`
		Color   string  `gorm:"column:color"`
		Total   float64 `gorm:"column:total"`
	}

	err := p.gorm.Table("expenses").
		Select("EXTRACT(YEAR FROM expanse_date) AS year, TO_CHAR(expanse_date, 'Month') AS month, tags.name AS tag_name, tags.color AS color, SUM(expenses.amount) AS total").
		Joins("INNER JOIN expense_tags ON expenses.id = expense_tags.expenses_id").
		Joins("INNER JOIN tags ON expense_tags.tags_id = tags.id").
		Where("expenses.user_id = ? AND EXTRACT(YEAR FROM expenses.expanse_date) = ? AND expenses.active = ?", userID, year, true).
		Group("year, month, tags.name, tags.color").
		Order("MIN(expanse_date)").
		Scan(&results).Error

	if err != nil {
		return nil, []int{}, errors.New("failed to fetch monthly expenses by tag: " + err.Error())
	}

	var years []int

	err = p.gorm.Table("expenses").
		Select("DISTINCT EXTRACT(YEAR FROM expanse_date) AS year").
		Where("expenses.user_id = ? AND expenses.active = ?", userID, true).
		Order("year").
		Scan(&years).Error

	if err != nil {
		return nil, nil, errors.New("failed to fetch available years: " + err.Error())
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

	var total float64
	var month string

	now := time.Now().In(location)
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, location)
	endOfMonth := now

	if err := p.gorm.Table("expenses").
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ? AND expanse_date BETWEEN ? AND ? AND active = ?", userID, startOfMonth, endOfMonth, true).
		Scan(&total).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			total = 0
		} else {
			return 0, "", errors.New("failed to fetch total expenses: " + err.Error())
		}
	}

	month = time.Now().Format("January")

	return total, month, nil
}

func (p *PresentersRepository) GetExpensesByMonthYear(userID string, month int, year int) (repositories.MonthExpenses, error) {
	var monthExpenses repositories.MonthExpenses
	monthExpenses.Month = time.Month(month).String()
	monthExpenses.Year = year

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Nanosecond)

	var expenses []Expenses
	if err := p.gorm.Where("user_id = ? AND expanse_date BETWEEN ? AND ? AND active = ?", userID, startDate, endDate, true).
		Find(&expenses).Error; err != nil {
		return repositories.MonthExpenses{}, errors.New("failed to fetch expenses: " + err.Error())
	}

	weeks := make(map[int]map[string]*repositories.DayExpense)
	totalExpenses := 0.0

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
		if err := p.gorm.Table("tags").
			Joins("JOIN expense_tags ON tags.id = expense_tags.tags_id").
			Where("expense_tags.expenses_id = ?", expense.ID).
			Select("tags.name, tags.color").
			Find(&tags).Error; err != nil {
			return repositories.MonthExpenses{}, errors.New("failed to fetch tags for expense: " + err.Error())
		}

		weeks[weekNumber][dayKey].Total += expense.Amount
		totalExpenses += expense.Amount

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

	monthExpenses.TotalExpenses = totalExpenses

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

	var availableYears []int
	if err := p.gorm.Table("expenses").
		Select("DISTINCT EXTRACT(YEAR FROM expanse_date) as year").
		Where("user_id = ? AND active = ?", userID, true).
		Order("year DESC").
		Pluck("year", &availableYears).Error; err != nil {
		return repositories.MonthExpenses{}, errors.New("failed to fetch available years: " + err.Error())
	}

	monthExpenses.AvailableYears = availableYears

	return monthExpenses, nil
}

func (p *PresentersRepository) GetTotalExpensesForCurrentWeek(userID string) (float64, string, error) {
	location, err := time.LoadLocation(util.TIMEZONE)
	if err != nil {
		return 0, "", errors.New("failed to load timezone: " + err.Error())
	}

	var totalExpenses float64

	now := time.Now().In(location)

	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	endOfWeek := now

	var expenses []Expenses
	if err := p.gorm.
		Where("user_id = ? AND expanse_date BETWEEN ? AND ?", userID, startOfMonth, endOfWeek).
		Find(&expenses).Error; err != nil {
		return 0, "", errors.New("failed to fetch expenses: " + err.Error())
	}

	for _, expense := range expenses {
		totalExpenses += expense.Amount
	}

	weekInterval := fmt.Sprintf("%s - %s", startOfMonth.Format("02/01/2006"), endOfWeek.Format("02/01/2006"))

	return totalExpenses, weekInterval, nil
}

func (p *PresentersRepository) GetTotalExpensesMonthCurrentYear(userID string, year int) (repositories.ExpensesMonthCurrentYear, error) {
	var expensesMonthCurrentYear repositories.ExpensesMonthCurrentYear
	expensesMonthCurrentYear.Year = year

	months := make([]repositories.MonthCurrentYear, 12)
	for i := 0; i < 12; i++ {
		months[i] = repositories.MonthCurrentYear{
			Month: time.Month(i + 1).String(),
			Total: 0,
		}
	}

	type ExpenseMonth struct {
		Month int     `json:"month"`
		Total float64 `json:"total"`
	}

	var expenses []ExpenseMonth
	if err := p.gorm.Table("expenses").
		Select("EXTRACT(MONTH FROM expanse_date) as month, COALESCE(SUM(amount), 0) as total").
		Where("user_id = ? AND EXTRACT(YEAR FROM expanse_date) = ? AND active = ?", userID, year, true).
		Group("EXTRACT(MONTH FROM expanse_date)").
		Order("EXTRACT(MONTH FROM expanse_date)").
		Find(&expenses).Error; err != nil {
		return repositories.ExpensesMonthCurrentYear{}, errors.New("failed to fetch expenses by month: " + err.Error())
	}

	for _, expense := range expenses {
		months[expense.Month-1].Total = expense.Total
	}

	var totalYear float64
	for _, month := range months {
		totalYear += month.Total
	}
	expensesMonthCurrentYear.Total = totalYear
	expensesMonthCurrentYear.Months = months

	var availableYears []int
	if err := p.gorm.Table("expenses").
		Select("DISTINCT EXTRACT(YEAR FROM expanse_date) as year").
		Where("user_id = ? AND active = ?", userID, true).
		Order("year DESC").
		Pluck("year", &availableYears).Error; err != nil {
		return repositories.ExpensesMonthCurrentYear{}, errors.New("failed to fetch available years: " + err.Error())
	}
	expensesMonthCurrentYear.AvailableYears = availableYears

	return expensesMonthCurrentYear, nil
}

func (p *PresentersRepository) GetCategoryTagsTotalsByMonthYear(userID string, month int, year int) (repositories.CategoryTagsTotals, error) {
	var categoryTagsTotals repositories.CategoryTagsTotals
	categoryTagsTotals.Month = time.Month(month).String()
	categoryTagsTotals.Year = year

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Nanosecond)

	var totalExpenses float64
	if err := p.gorm.Table("expenses").
		Where("user_id = ? AND expanse_date BETWEEN ? AND ? AND active = ?", userID, startDate, endDate, true).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalExpenses).Error; err != nil {
		return repositories.CategoryTagsTotals{}, errors.New("failed to calculate total expenses for the month: " + err.Error())
	}
	categoryTagsTotals.ExpensesAmount = totalExpenses

	var results []struct {
		CategoryName  string
		CategoryTotal float64
	}

	if err := p.gorm.Table("expenses").
		Select("categories.name as category_name, COALESCE(SUM(expenses.amount), 0) as category_total").
		Joins("LEFT JOIN categories ON categories.id = expenses.category_id").
		Where("expenses.user_id = ? AND expenses.expanse_date BETWEEN ? AND ? AND expenses.active = ?", userID, startDate, endDate, true).
		Group("categories.name").
		Scan(&results).Error; err != nil {
		return repositories.CategoryTagsTotals{}, errors.New("failed to fetch expenses by category: " + err.Error())
	}

	categoryMap := make(map[string]*repositories.CategoryWithTags)
	for _, result := range results {
		categoryMap[result.CategoryName] = &repositories.CategoryWithTags{
			Name:           result.CategoryName,
			CategoryAmount: result.CategoryTotal,
			Tags:           []repositories.CategoryTagTotal{},
		}
	}

	var resultsTags []struct {
		CategoryName string
		TagName      string
		TagTotal     float64
	}

	if err := p.gorm.Table("expenses").
		Select("categories.name as category_name, tags.name as tag_name, COALESCE(SUM(expenses.amount), 0) as tag_total").
		Joins("LEFT JOIN categories ON categories.id = expenses.category_id").
		Joins("LEFT JOIN expense_tags ON expense_tags.expenses_id = expenses.id").
		Joins("LEFT JOIN tags ON tags.id = expense_tags.tags_id").
		Where("expenses.user_id = ? AND expenses.expanse_date BETWEEN ? AND ? AND expenses.active = ?", userID, startDate, endDate, true).
		Group("categories.name, tags.name").
		Scan(&resultsTags).Error; err != nil {
		return repositories.CategoryTagsTotals{}, errors.New("failed to fetch expenses by category and tags: " + err.Error())
	}

	for _, result := range resultsTags {
		if category, exists := categoryMap[result.CategoryName]; exists {
			category.Tags = append(category.Tags, repositories.CategoryTagTotal{
				Name:      result.TagName,
				TagAmount: result.TagTotal,
			})
		}
	}

	for _, category := range categoryMap {
		categoryTagsTotals.Categories = append(categoryTagsTotals.Categories, *category)
	}

	sort.Slice(categoryTagsTotals.Categories, func(i, j int) bool {
		return categoryTagsTotals.Categories[i].Name < categoryTagsTotals.Categories[j].Name
	})

	var availableYears []int
	if err := p.gorm.Table("expenses").
		Distinct("EXTRACT(YEAR FROM expanse_date)").
		Where("user_id = ? AND active = ?", userID, true).
		Order("EXTRACT(YEAR FROM expanse_date) DESC").
		Pluck("EXTRACT(YEAR FROM expanse_date)", &availableYears).Error; err != nil {
		return repositories.CategoryTagsTotals{}, errors.New("failed to fetch available years: " + err.Error())
	}
	categoryTagsTotals.AvailableYears = availableYears

	var availableMonths []struct {
		Month int
	}
	if err := p.gorm.Table("expenses").
		Select("DISTINCT EXTRACT(MONTH FROM expanse_date) AS month").
		Where("user_id = ? AND EXTRACT(YEAR FROM expanse_date) = ? AND active = ?", userID, year, true).
		Order("month ASC").
		Scan(&availableMonths).Error; err != nil {
		return repositories.CategoryTagsTotals{}, errors.New("failed to fetch available months: " + err.Error())
	}

	for _, m := range availableMonths {
		categoryTagsTotals.AvailableMonths = append(categoryTagsTotals.AvailableMonths, repositories.MonthOption{
			Label: time.Month(m.Month).String(),
			Value: fmt.Sprintf("%02d", m.Month),
		})
	}

	return categoryTagsTotals, nil
}

func (p *PresentersRepository) GetAvailableMonthsYears(userID string) ([]int, []repositories.MonthOption, error) {
	var availableYears []int
	if err := p.gorm.Table("expenses").
		Select("DISTINCT EXTRACT(YEAR FROM expanse_date) as year").
		Where("user_id = ? AND active = ?", userID, true).
		Order("year DESC").
		Pluck("year", &availableYears).Error; err != nil {
		return nil, nil, errors.New("failed to fetch available years: " + err.Error())
	}

	var monthOptions []repositories.MonthOption
	for i := 1; i <= 12; i++ {
		monthOptions = append(monthOptions, repositories.MonthOption{
			Label: time.Month(i).String(),
			Value: fmt.Sprintf("%02d", i),
		})
	}

	return availableYears, monthOptions, nil
}

func (p *PresentersRepository) GetDayToDayExpensesPeriod(userID string, startDate time.Time, endDate time.Time) ([]entities.Expense, error) {
	var expensesModel []Expenses

	if err := p.gorm.
		Where("user_id = ? AND active = ? AND expanse_date BETWEEN ? AND ?", userID, true, startDate, endDate).
		Find(&expensesModel).Error; err != nil {
		return []entities.Expense{}, errors.New("failed to fetch expenses: " + err.Error())
	}

	var expenses []entities.Expense

	if len(expensesModel) > 0 {
		for _, expenseModel := range expensesModel {

			expense := entities.Expense{
				SharedEntity: entities.SharedEntity{
					ID: expenseModel.ID,
				},
				Amount:      expenseModel.Amount,
				ExpenseDate: expenseModel.ExpanseDate,
			}

			expenses = append(expenses, expense)
		}

		sort.Slice(expenses, func(i, j int) bool {
			return expenses[i].ExpenseDate.After(expenses[j].ExpenseDate)
		})
	} else {
		expenses = []entities.Expense{}
	}

	return expenses, nil
}
