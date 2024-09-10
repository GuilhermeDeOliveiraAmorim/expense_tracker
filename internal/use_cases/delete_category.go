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
	user, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return DeleteCategoryOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   err.Error(),
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
