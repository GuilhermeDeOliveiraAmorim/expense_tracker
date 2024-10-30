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
	SuccessMessage string `json:"success_message"`
	ContentMessage string `json:"content_message"`
}

type DeleteCategoryUseCase struct {
	CategoryRepository repositories.CategoryRepositoryInterface
	UserRepository     repositories.UserRepositoryInterface
}

func NewDeleteCategoryUseCase(
	CategoryRepository repositories.CategoryRepositoryInterface,
	UserRepository repositories.UserRepositoryInterface,
) *DeleteCategoryUseCase {
	return &DeleteCategoryUseCase{
		CategoryRepository: CategoryRepository,
		UserRepository:     UserRepository,
	}
}

func (c *DeleteCategoryUseCase) Execute(input DeleteCategoryInputDto) (DeleteCategoryOutputDto, []util.ProblemDetails) {
	user, GetUserErr := c.UserRepository.GetUser(input.UserID)
	if GetUserErr != nil {
		return DeleteCategoryOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   GetUserErr.Error(),
				Instance: util.RFC404,
			},
		}
	} else if !user.Active {
		return DeleteCategoryOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Forbidden",
				Title:    "User is not active",
				Status:   403,
				Detail:   "User is not active",
				Instance: util.RFC403,
			},
		}
	}

	categoryToDelete, GetCategoryErr := c.CategoryRepository.GetCategory(input.UserID, input.CategoryID)
	if GetCategoryErr != nil {
		return DeleteCategoryOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "Category not found",
				Status:   404,
				Detail:   GetCategoryErr.Error(),
				Instance: util.RFC404,
			},
		}
	}

	categoryToDelete.Deactivate()

	deleteCategoryErr := c.CategoryRepository.DeleteCategory(categoryToDelete)
	if deleteCategoryErr != nil {
		if deleteCategoryErr.Error() == "there are expenses associated with this category" {
			return DeleteCategoryOutputDto{}, []util.ProblemDetails{
				{
					Type:     "Conflict",
					Title:    "Category has expenses",
					Status:   409,
					Detail:   "Error: " + deleteCategoryErr.Error(),
					Instance: util.RFC409,
				},
			}
		}

		return DeleteCategoryOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Err deleting category",
				Status:   500,
				Detail:   deleteCategoryErr.Error(),
				Instance: util.RFC500,
			},
		}
	}

	return DeleteCategoryOutputDto{
		SuccessMessage: "Category deleted successfully",
		ContentMessage: "Category " + categoryToDelete.Name + " deleted",
	}, nil
}
