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
}

func NewPresentersFactory(db *gorm.DB) *PresentersFactory {
	presentersRepository := repositoriesgorm.NewPresentersRepository(db)
	userRepository := repositoriesgorm.NewUserRepository(db)

	getTotalExpensesForPeriod := presenters.NewGetTotalExpensesForPeriodUseCase(presentersRepository, userRepository)
	getExpensesByCategoryPeriod := presenters.NewGetExpensesByCategoryPeriodUseCase(presentersRepository, userRepository)
	getMonthlyExpensesByCategoryYear := presenters.NewGetMonthlyExpensesByCategoryYearUseCase(presentersRepository, userRepository)
	getMonthlyExpensesByTagYear := presenters.NewGetMonthlyExpensesByTagYearUseCase(presentersRepository, userRepository)

	return &PresentersFactory{
		GetTotalExpensesForPeriod:        getTotalExpensesForPeriod,
		GetExpensesByCategoryPeriod:      getExpensesByCategoryPeriod,
		GetMonthlyExpensesByCategoryYear: getMonthlyExpensesByCategoryYear,
		GetMonthlyExpensesByTagYear:      getMonthlyExpensesByTagYear,
	}
}
