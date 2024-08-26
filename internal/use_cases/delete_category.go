package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
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

func NewDeleteCategoryUseCase(
	CategoryRepository repositories.CategoryRepositoryInterface,
) *DeleteCategoryUseCase {
	return &DeleteCategoryUseCase{
		CategoryRepository: CategoryRepository,
	}
}

func (c *DeleteCategoryUseCase) Execute(input DeleteCategoryInputDto) (DeleteCategoryOutputDto, []util.ProblemDetails) {
	categoryToDelete, err := c.CategoryRepository.GetCategory(input.CategoryID)
	if err != nil {
		return DeleteCategoryOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "Category not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC400,
			},
		}
	}

	categoryToDelete.Deactivate()

	err = c.CategoryRepository.DeleteCategory(categoryToDelete)
	if err != nil {
		return DeleteCategoryOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Err deleting category",
				Status:   500,
				Detail:   err.Error(),
				Instance: util.RFC500,
			},
		}
	}

	return DeleteCategoryOutputDto{
		ID: categoryToDelete.ID,
	}, nil
}
