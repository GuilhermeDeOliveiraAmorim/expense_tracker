package factory

import (
	repositoriesgorm "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/infra/repositories_gorm"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/presenters"
	"gorm.io/gorm"
)

type PresentersFactory struct {
	GetTotalExpensesForPeriod *presenters.GetTotalExpensesForPeriodUseCase
}

func NewPresentersFactory(db *gorm.DB) *PresentersFactory {
	presentersRepository := repositoriesgorm.NewPresentersRepository(db)
	userRepository := repositoriesgorm.NewUserRepository(db)

	getTotalExpensesForPeriod := presenters.NewGetTotalExpensesForPeriodUseCase(presentersRepository, userRepository)

	return &PresentersFactory{
		GetTotalExpensesForPeriod: getTotalExpensesForPeriod,
	}
}
