package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type GetCategoriesInputDto struct {
	UserID string `json:"user_id"`
}

type GetCategoriesOutputDto struct {
	Categories []entities.Category `json:"categories"`
}

type GetCategoriesUseCase struct {
	CategoryRepository repositories.CategoryRepositoryInterface
}

func NewGetCategoriesUseCase(
	CategoryRepository repositories.CategoryRepositoryInterface,
) *GetCategoriesUseCase {
	return &GetCategoriesUseCase{
		CategoryRepository: CategoryRepository,
	}
}

func (c *GetCategoriesUseCase) Execute(input GetCategoriesInputDto) (GetCategoriesOutputDto, []util.ProblemDetails) {
	searchedsCategories, err := c.CategoryRepository.GetCategories(input.UserID)
	if err != nil {
		return GetCategoriesOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Err fetching categories",
				Status:   500,
				Detail:   err.Error(),
				Instance: util.RFC500,
			},
		}
	}

	return GetCategoriesOutputDto{
		Categories: searchedsCategories,
	}, nil
}
