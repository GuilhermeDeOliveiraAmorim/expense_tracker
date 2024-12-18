package factory

import (
	repositoriesgorm "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/infra/repositories_gorm"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/presenters"
	"gorm.io/gorm"
)

type PresentersFactory struct {
	GetTotalExpensesForPeriod        *presenters.GetTotalExpensesForPeriodUseCase
	GetExpensesByCategoryPeriod      *presenters.GetExpensesByCategoryPeriodUseCase
	GetMonthlyExpensesByCategoryYear *presenters.GetMonthlyExpensesByCategoryYearUseCase
	GetMonthlyExpensesByTagYear      *presenters.GetMonthlyExpensesByTagYearUseCase
	GetTotalExpensesForCurrentMonth  *presenters.GetTotalExpensesForCurrentMonthUseCase
	GetTotalExpensesForCurrentWeek   *presenters.GetTotalExpensesForCurrentWeekUseCase
	GetExpensesByMonthYear           *presenters.GetExpensesByMonthYearUseCase
	GetTotalExpensesMonthCurrentYear *presenters.GetTotalExpensesMonthCurrentYearUseCase
	GetCategoryTagsTotalsByMonthYear *presenters.GetCategoryTagsTotalsByMonthYearUseCase
	GetAvailableMonthsYears          *presenters.GetAvailableMonthsYearsUseCase
	GetDayToDayExpensesPeriod        *presenters.GetDayToDayExpensesPeriodUseCase
}

func NewPresentersFactory(db *gorm.DB) *PresentersFactory {
	presentersRepository := repositoriesgorm.NewPresentersRepository(db)
	userRepository := repositoriesgorm.NewUserRepository(db)

	getTotalExpensesForPeriod := presenters.NewGetTotalExpensesForPeriodUseCase(presentersRepository, userRepository)
	getExpensesByCategoryPeriod := presenters.NewGetExpensesByCategoryPeriodUseCase(presentersRepository, userRepository)
	getMonthlyExpensesByCategoryYear := presenters.NewGetMonthlyExpensesByCategoryYearUseCase(presentersRepository, userRepository)
	getMonthlyExpensesByTagYear := presenters.NewGetMonthlyExpensesByTagYearUseCase(presentersRepository, userRepository)
	getTotalExpensesForCurrentMonth := presenters.NewGetTotalExpensesForCurrentMonthUseCase(presentersRepository, userRepository)
	getTotalExpensesForCurrentWeek := presenters.NewGetTotalExpensesForCurrentWeekUseCase(presentersRepository, userRepository)
	getExpensesByMonthYear := presenters.NewGetExpensesByMonthYearUseCase(presentersRepository, userRepository)
	getTotalExpensesMonthCurrentYear := presenters.NewGetTotalExpensesMonthCurrentYearUseCase(presentersRepository, userRepository)
	getCategoryTagsTotalsByMonthYear := presenters.NewGetCategoryTagsTotalsByMonthYearUseCase(presentersRepository, userRepository)
	getAvailableMonthsYears := presenters.NewGetAvailableMonthsYearsUseCase(presentersRepository, userRepository)
	getDayToDayExpensesPeriod := presenters.NewGetDayToDayExpensesPeriodUseCase(presentersRepository, userRepository)

	return &PresentersFactory{
		GetTotalExpensesForPeriod:        getTotalExpensesForPeriod,
		GetExpensesByCategoryPeriod:      getExpensesByCategoryPeriod,
		GetMonthlyExpensesByCategoryYear: getMonthlyExpensesByCategoryYear,
		GetMonthlyExpensesByTagYear:      getMonthlyExpensesByTagYear,
		GetTotalExpensesForCurrentMonth:  getTotalExpensesForCurrentMonth,
		GetTotalExpensesForCurrentWeek:   getTotalExpensesForCurrentWeek,
		GetExpensesByMonthYear:           getExpensesByMonthYear,
		GetTotalExpensesMonthCurrentYear: getTotalExpensesMonthCurrentYear,
		GetCategoryTagsTotalsByMonthYear: getCategoryTagsTotalsByMonthYear,
		GetAvailableMonthsYears:          getAvailableMonthsYears,
		GetDayToDayExpensesPeriod:        getDayToDayExpensesPeriod,
	}
}
