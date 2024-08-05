package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
)

type GetCategoriesInputDto struct {
}

type GetCategoriesOutputDto struct {
	Categories []entities.Category `json:"categories"`
}

type GetCategoriesUseCase struct {
	CategoryRepository repositories.CategoryRepositoryInterface
}

func (c *GetCategoriesUseCase) Execute(input GetCategoriesInputDto) (GetCategoriesOutputDto, []error) {
	searchedsCategories, err := c.CategoryRepository.GetCategories()
	if err != nil {
		return GetCategoriesOutputDto{}, err
	}

	return GetCategoriesOutputDto{
		Categories: searchedsCategories,
	}, nil
}
