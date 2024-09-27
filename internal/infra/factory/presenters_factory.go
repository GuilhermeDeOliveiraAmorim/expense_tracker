package factory

import (
	"gorm.io/gorm"
)

type PresentersFactory struct {
}

func NewPresentersFactory(db *gorm.DB) *PresentersFactory {
	//presentersRepository := repositoriesgorm.NewPresentersRepository(db)
	//userRepository := repositoriesgorm.NewUserRepository(db)

	return &PresentersFactory{}
}
