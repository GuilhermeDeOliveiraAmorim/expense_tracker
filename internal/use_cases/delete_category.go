package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type DeleteCategoryInputDto struct {
	UserID     string `json:"user_id"`
	CategoryID string `json:"category_id"`
}

type DeleteCategoryOutputDto struct {
	Message string `json:"message"`
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
	categoryToDelete, err := c.CategoryRepository.GetCategory(input.UserID, input.CategoryID)
	if err != nil {
		return DeleteCategoryOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "Category not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
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
		Message: "Category deleted successfully",
	}, nil
}
