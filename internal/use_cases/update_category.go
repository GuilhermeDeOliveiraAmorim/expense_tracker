package usecases

import (
	"fmt"
	"strings"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type UpdateCategoryInputDto struct {
	CategoryID string `json:"category_id"`
	Name       string `json:"name"`
}

type UpdateCategoryOutputDto struct {
	CategoryID string `json:"category_id"`
}

type UpdateCategoryUseCase struct {
	CategoryRepository repositories.CategoryRepositoryInterface
}

func NewUpdateCategoryUseCase(
	CategoryRepository repositories.CategoryRepositoryInterface,
) *UpdateCategoryUseCase {
	return &UpdateCategoryUseCase{
		CategoryRepository: CategoryRepository,
	}
}

func (c *UpdateCategoryUseCase) Execute(input UpdateCategoryInputDto) (UpdateCategoryOutputDto, []util.ProblemDetails) {
	existingCategory, GetCategoryByNameErr := c.CategoryRepository.ThisCategoryExists(input.Name)
	if GetCategoryByNameErr != nil && strings.Compare(GetCategoryByNameErr.Error(), "category not found") > 0 {
		return UpdateCategoryOutputDto{}, []util.ProblemDetails{
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
		return UpdateCategoryOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Validation Error",
				Title:    "Category already exists",
				Status:   409,
				Detail:   "A category with this name already exists",
				Instance: util.RFC409,
			},
		}
	}

	searchedCategory, err := c.CategoryRepository.GetCategory(input.CategoryID)
	if err != nil {
		return UpdateCategoryOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "Category not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC400,
			},
		}
	}

	if input.Name != searchedCategory.Name {
		err := searchedCategory.ChangeName(input.Name)
		if len(err) > 0 {
			return UpdateCategoryOutputDto{}, err
		}

		UpdateCategoryErr := c.CategoryRepository.UpdateCategory(searchedCategory)
		if err != nil {
			return UpdateCategoryOutputDto{}, []util.ProblemDetails{
				{
					Type:     "Validation Error",
					Title:    "Bad Request",
					Status:   500,
					Detail:   UpdateCategoryErr.Error(),
					Instance: util.RFC500,
				},
			}
		}
	} else {
		return UpdateCategoryOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Validation Error",
				Title:    "Bad Request",
				Status:   400,
				Detail:   fmt.Sprintf("Category name is already '%s'", searchedCategory.Name),
				Instance: util.RFC400,
			},
		}
	}

	return UpdateCategoryOutputDto{
		CategoryID: searchedCategory.ID,
	}, nil
}
