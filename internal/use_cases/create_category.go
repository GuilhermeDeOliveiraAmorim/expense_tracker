package usecases

import (
	"strings"

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
	existingCategory, GetCategoryByNameErr := c.CategoryRepository.ThisCategoryExists(input.Name)
	if GetCategoryByNameErr != nil && strings.Compare(GetCategoryByNameErr.Error(), "category not found") > 0 {
		return CreateCategoryOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Error fetching existing category",
				Status:   500,
				Detail:   GetCategoryByNameErr.Error(),
				Instance: util.RFC500,
			},
		}
	}

	if existingCategory {
		return CreateCategoryOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Validation Error",
				Title:    "Category already exists",
				Status:   409,
				Detail:   "A category with this name already exists",
				Instance: util.RFC409,
			},
		}
	}

	newCategory, err := entities.NewCategory(input.Name)
	if err != nil {
		return CreateCategoryOutputDto{}, err
	}

	CreateCategoryErr := c.CategoryRepository.CreateCategory(*newCategory)
	if CreateCategoryErr != nil {
		return CreateCategoryOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Error creating new category",
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
