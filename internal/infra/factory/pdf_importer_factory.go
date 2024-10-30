package factory

import (
	pdfservice "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/infra/pdf_service"
	repositoriesgorm "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/infra/repositories_gorm"
	usecases "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/use_cases"
	"gorm.io/gorm"
)

type PdfImporterFactory struct {
	PdfImporter *usecases.PdfImporterUseCase
}

func NewPdfImporterFactory(db *gorm.DB) *PdfImporterFactory {
	pdfImporterRepository := pdfservice.NewPdfImporterRepository(db)
	userRepository := repositoriesgorm.NewUserRepository(db)

	pdfImporter := usecases.NewPdfImporterUseCase(userRepository, pdfImporterRepository)

	return &PdfImporterFactory{
		PdfImporter: pdfImporter,
	}
}
