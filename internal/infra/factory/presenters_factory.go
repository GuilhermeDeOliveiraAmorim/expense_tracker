package factory

import (
	repositoriesgorm "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/infra/repositories_gorm"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/presenters"
	"gorm.io/gorm"
)

type PresentersFactory struct {
	ShowTotalExpensesCategoryPeriod *presenters.ShowTotalExpensesCategoryPeriodPresenters
	ShowCategoryTreemapAmountPeriod *presenters.ShowCategoryTreemapAmountPeriodPresenters
	ShowExpenseSimpleTablePeriod    *presenters.ShowExpenseSimpleTablePeriodPresenters
}

func NewPresentersFactory(db *gorm.DB) *PresentersFactory {
	presentersRepository := repositoriesgorm.NewPresentersRepository(db)
	userRepository := repositoriesgorm.NewUserRepository(db)

	showTotalExpensesCategoryPeriod := presenters.NewShowTotalExpensesCategoryPeriodPresenters(presentersRepository, userRepository)
	showCategoryTreemapAmountPeriod := presenters.NewCategoryTreemapAmountPeriodPresenters(presentersRepository, userRepository)
	showExpenseSimpleTablePeriod := presenters.NewShowExpenseSimpleTablePeriodPresenters(presentersRepository, userRepository)

	return &PresentersFactory{
		ShowTotalExpensesCategoryPeriod: showTotalExpensesCategoryPeriod,
		ShowCategoryTreemapAmountPeriod: showCategoryTreemapAmountPeriod,
		ShowExpenseSimpleTablePeriod:    showExpenseSimpleTablePeriod,
	}
}
