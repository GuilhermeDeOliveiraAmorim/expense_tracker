package usecases

import (
	"fmt"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
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

func NewUpdateCategoryUseCase(
	CategoryRepository repositories.CategoryRepositoryInterface,
) *UpdateCategoryUseCase {
	return &UpdateCategoryUseCase{
		CategoryRepository: CategoryRepository,
	}
}

func (c *UpdateCategoryUseCase) Execute(input UpdateCategoryInputDto) (UpdateCategoryOutputDto, []util.ProblemDetails) {
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
		ID: searchedCategory.ID,
	}, nil
}
