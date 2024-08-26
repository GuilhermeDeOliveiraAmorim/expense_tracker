package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
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

func NewCreateCategoryUseCase(
	CategoryRepository repositories.CategoryRepositoryInterface,
) *CreateCategoryUseCase {
	return &CreateCategoryUseCase{
		CategoryRepository: CategoryRepository,
	}
}

func (c *CreateCategoryUseCase) Execute(input CreateCategoryInputDto) (CreateCategoryOutputDto, []util.ProblemDetails) {
	newCategory, err := entities.NewCategory(input.Name)
	if err != nil {
		return CreateCategoryOutputDto{}, err
	}

	CreateCategoryErr := c.CategoryRepository.CreateCategory(*newCategory)
	if err != nil {
		return CreateCategoryOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Validation Error",
				Title:    "Bad Request",
				Status:   500,
				Detail:   CreateCategoryErr.Error(),
				Instance: util.RFC500,
			},
		}
	}

	return CreateCategoryOutputDto{
		ID: newCategory.ID,
	}, nil
}
