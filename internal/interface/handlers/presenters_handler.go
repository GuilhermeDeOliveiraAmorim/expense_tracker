package handlers

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/infra/factory"
)

type PresentersHandler struct {
	presenterFactory *factory.PresentersFactory
}

func NewPresentersHandler(factory *factory.PresentersFactory) *PresentersHandler {
	return &PresentersHandler{
		presenterFactory: factory,
	}
}
