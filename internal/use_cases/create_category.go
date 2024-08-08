package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
)

type CreateCategoryInputDto struct {
	Name string `json:"name"`
}

type CreateCategoryOutputDto struct {
	ID string `json:"id"`
}

type CreateCategoryUseCase struct {
	CategoryRepository repositories.CategoryRepositoryInterface
}

func (c *CreateCategoryUseCase) Execute(input CreateCategoryInputDto) (CreateCategoryOutputDto, []error) {
	newCategory, err := entities.NewCategory(input.Name)
	if err != nil {
		return CreateCategoryOutputDto{}, err
	}

	err = c.CategoryRepository.CreateCategory(*newCategory)
	if err != nil {
		return CreateCategoryOutputDto{}, err
	}

	return CreateCategoryOutputDto{
		ID: newCategory.ID,
	}, nil
}
