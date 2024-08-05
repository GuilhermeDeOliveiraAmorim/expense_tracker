package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
)

type DeleteCategoryInputDto struct {
	CategoryID string `json:"id"`
}

type DeleteCategoryOutputDto struct {
	ID string `json:"id"`
}

type DeleteCategoryUseCase struct {
	CategoryRepository repositories.CategoryRepositoryInterface
}

func (c *DeleteCategoryUseCase) Execute(input DeleteCategoryInputDto) (DeleteCategoryOutputDto, []error) {
	categoryToDelete, err := c.CategoryRepository.GetCategory(input.CategoryID)
	if err != nil {
		return DeleteCategoryOutputDto{}, err
	}

	categoryToDelete.Deactivate()

	err = c.CategoryRepository.DeleteCategory(categoryToDelete)
	if err != nil {
		return DeleteCategoryOutputDto{}, err
	}

	return DeleteCategoryOutputDto{
		ID: categoryToDelete.ID,
	}, nil
}
