package factory

import (
	repositoriesgorm "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/infra/repositories_gorm"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/presenters"
	"gorm.io/gorm"
)

type PresentersFactory struct {
	ShowTotalExpensesCategoryPeriod *presenters.ShowTotalExpensesCategoryPeriodPresenters
	ShowCategoryTreemapAmountPeriod *presenters.ShowCategoryTreemapAmountPeriodPresenters
}

func NewPresentersFactory(db *gorm.DB) *PresentersFactory {
	userRepository := repositoriesgorm.NewUserRepository(db)
	presentersRepository := repositoriesgorm.NewPresentersRepository(db)

	showTotalExpensesCategoryPeriod := presenters.NewShowTotalExpensesCategoryPeriodPresenters(presentersRepository, userRepository)
	showCategoryTreemapAmountPeriod := presenters.NewCategoryTreemapAmountPeriodPresenters(presentersRepository, userRepository)

	return &PresentersFactory{
		ShowTotalExpensesCategoryPeriod: showTotalExpensesCategoryPeriod,
		ShowCategoryTreemapAmountPeriod: showCategoryTreemapAmountPeriod,
	}
}
