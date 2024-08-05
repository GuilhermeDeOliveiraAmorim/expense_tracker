package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
)

type AddCategoryInputDto struct {
	Name string `json:"name"`
}

type AddCategoryOutputDto struct {
	ID string `json:"id"`
}

type AddCategoryUseCase struct {
	CategoryRepository repositories.CategoryRepositoryInterface
}

func (c *AddCategoryUseCase) Execute(input AddCategoryInputDto) (AddCategoryOutputDto, []error) {
	newCategory, err := entities.NewCategory(input.Name)
	if err != nil {
		return AddCategoryOutputDto{}, err
	}

	err = c.CategoryRepository.CreateCategory(*newCategory)
	if err != nil {
		return AddCategoryOutputDto{}, err
	}

	return AddCategoryOutputDto{
		ID: newCategory.ID,
	}, nil
}
