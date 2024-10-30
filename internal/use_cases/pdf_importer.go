package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type PdfImporterInputDto struct {
	UserID            string `json:"user_id"`
	DefaultCategoryID string `json:"default_category_id"`
	Path              string `json:"path"`
}

type PdfImporterOutputDto struct {
	SuccessMessage string `json:"success_message"`
	ContentMessage string `json:"content_message"`
}

type PdfImporterUseCase struct {
	UserRepository        repositories.UserRepositoryInterface
	PdfImporterRepository repositories.PdfImporterRepositoryInterface
}

func NewPdfImporterUseCase(
	UserRepository repositories.UserRepositoryInterface,
	PdfImporterRepository repositories.PdfImporterRepositoryInterface,
) *PdfImporterUseCase {
	return &PdfImporterUseCase{
		UserRepository:        UserRepository,
		PdfImporterRepository: PdfImporterRepository,
	}
}

func (c *PdfImporterUseCase) Execute(input PdfImporterInputDto) (PdfImporterOutputDto, []util.ProblemDetails) {
	user, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return PdfImporterOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
			},
		}
	} else if !user.Active {
		return PdfImporterOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Forbidden",
				Title:    "User is not active",
				Status:   403,
				Detail:   "User is not active",
				Instance: util.RFC403,
			},
		}
	}

	errImportPdf := c.PdfImporterRepository.ImportPdf(input.UserID, input.DefaultCategoryID, input.Path)
	if errImportPdf != nil {
		return PdfImporterOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Error importing PDF",
				Status:   500,
				Detail:   errImportPdf.Error(),
				Instance: util.RFC500,
			},
		}
	}

	return PdfImporterOutputDto{
		SuccessMessage: "PDF imported successfully",
		ContentMessage: "PDF content: " + input.Path,
	}, nil
}
