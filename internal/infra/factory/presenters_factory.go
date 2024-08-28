package factory

import (
	repositoriesgorm "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/infra/repositories_gorm"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/presenters"
	"gorm.io/gorm"
)

type PresentersFactory struct {
	ShowTotalExpensesCategoryPeriod *presenters.ShowTotalExpensesCategoryPeriodPresenters
}

func NewPresentersFactory(db *gorm.DB) *PresentersFactory {
	userRepository := repositoriesgorm.NewUserRepository(db)
	presentersRepository := repositoriesgorm.NewPresentersRepository(db)

	showTotalExpensesCategoryPeriod := presenters.NewShowTotalExpensesCategoryPeriodPresenters(presentersRepository, userRepository)

	return &PresentersFactory{
		ShowTotalExpensesCategoryPeriod: showTotalExpensesCategoryPeriod,
	}
}
