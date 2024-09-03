package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type UpdateCategoryInputDto struct {
	CategoryID string `json:"category_id"`
	Name       string `json:"name"`
	Color      string `json:"color"`
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
	searchedCategory, getCategoryErr := c.CategoryRepository.GetCategory(input.CategoryID)
	if getCategoryErr != nil {
		return UpdateCategoryOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "Category not found",
				Status:   404,
				Detail:   getCategoryErr.Error(),
				Instance: util.RFC400,
			},
		}
	}

	changeNameErr := searchedCategory.ChangeName(input.Name)
	if len(changeNameErr) > 0 {
		return UpdateCategoryOutputDto{}, changeNameErr
	}

	changeColorErr := searchedCategory.ChangeColor(input.Color)
	if len(changeColorErr) > 0 {
		return UpdateCategoryOutputDto{}, changeColorErr
	}

	updateCategoryErr := c.CategoryRepository.UpdateCategory(searchedCategory)
	if updateCategoryErr != nil {
		return UpdateCategoryOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Validation Error",
				Title:    "Bad Request",
				Status:   500,
				Detail:   updateCategoryErr.Error(),
				Instance: util.RFC500,
			},
		}
	}

	return UpdateCategoryOutputDto{
		CategoryID: searchedCategory.ID,
	}, nil
}
