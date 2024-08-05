package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
)

type GetCategoryInputDto struct {
	CategoryID string `json:"category_id"`
}

type GetCategoryOutputDto struct {
	Category entities.Category `json:"category"`
}

type GetCategoryUseCase struct {
	CategoryRepository repositories.CategoryRepositoryInterface
}

func (c *GetCategoryUseCase) Execute(input GetCategoryInputDto) (GetCategoryOutputDto, []error) {
	searchedCategory, err := c.CategoryRepository.GetCategory(input.CategoryID)
	if err != nil {
		return GetCategoryOutputDto{}, err
	}

	return GetCategoryOutputDto{
		Category: searchedCategory,
	}, nil
}
