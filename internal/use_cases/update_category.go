package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
)

type UpdateCategoryInputDto struct {
	CategoryID string `json:"id"`
	Name       string `json:"name"`
}

type UpdateCategoryOutputDto struct {
	ID string `json:"id"`
}

type UpdateCategoryUseCase struct {
	CategoryRepository repositories.CategoryRepositoryInterface
}

func (c *UpdateCategoryUseCase) Execute(input UpdateCategoryInputDto) (UpdateCategoryOutputDto, []error) {
	searchedCategory, err := c.CategoryRepository.GetCategory(input.CategoryID)
	if err != nil {
		return UpdateCategoryOutputDto{}, err
	}

	err = searchedCategory.ChangeName(input.Name)
	if len(err) > 0 {
		return UpdateCategoryOutputDto{}, err
	}

	err = c.CategoryRepository.UpdateCategory(searchedCategory)
	if err != nil {
		return UpdateCategoryOutputDto{}, err
	}

	return UpdateCategoryOutputDto{
		ID: searchedCategory.ID,
	}, nil
}
