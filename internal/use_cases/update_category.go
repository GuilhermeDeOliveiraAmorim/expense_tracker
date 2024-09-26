package usecases

import (
	"strings"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type UpdateCategoryInputDto struct {
	UserID     string `json:"user_id"`
	CategoryID string `json:"category_id"`
	Name       string `json:"name"`
	Color      string `json:"color"`
}

type UpdateCategoryOutputDto struct {
	CategoryID     string `json:"category_id"`
	SuccessMessage string `json:"success_message"`
	ContentMessage string `json:"content_message"`
}

type UpdateCategoryUseCase struct {
	CategoryRepository repositories.CategoryRepositoryInterface
	UserRepository     repositories.UserRepositoryInterface
}

func NewUpdateCategoryUseCase(
	CategoryRepository repositories.CategoryRepositoryInterface,
	UserRepository repositories.UserRepositoryInterface,
) *UpdateCategoryUseCase {
	return &UpdateCategoryUseCase{
		CategoryRepository: CategoryRepository,
		UserRepository:     UserRepository,
	}
}

func (c *UpdateCategoryUseCase) Execute(input UpdateCategoryInputDto) (UpdateCategoryOutputDto, []util.ProblemDetails) {
	user, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return UpdateCategoryOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
			},
		}
	} else if !user.Active {
		return UpdateCategoryOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Forbidden",
				Title:    "User is not active",
				Status:   403,
				Detail:   "User is not active",
				Instance: util.RFC403,
			},
		}
	}

	existingCategory, GetCategoryByNameErr := c.CategoryRepository.ThisCategoryExists(input.UserID, input.Name)
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

	searchedCategory, getCategoryErr := c.CategoryRepository.GetCategory(input.UserID, input.CategoryID)
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
		CategoryID:     searchedCategory.ID,
		SuccessMessage: "Category updated successfully",
		ContentMessage: "Category ID: " + searchedCategory.ID,
	}, nil
}
