package usecases

import (
	"strings"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type CreateCategoryInputDto struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
}

type CreateCategoryOutputDto struct {
	CategoryID string `json:"category_id"`
	Message    string `json:"message"`
}

type CreateCategoryUseCase struct {
	CategoryRepository repositories.CategoryRepositoryInterface
	UserRepository     repositories.UserRepositoryInterface
}

func NewCreateCategoryUseCase(
	CategoryRepository repositories.CategoryRepositoryInterface,
	UserRepository repositories.UserRepositoryInterface,
) *CreateCategoryUseCase {
	return &CreateCategoryUseCase{
		CategoryRepository: CategoryRepository,
		UserRepository:     UserRepository,
	}
}

func (c *CreateCategoryUseCase) Execute(input CreateCategoryInputDto) (CreateCategoryOutputDto, []util.ProblemDetails) {
	user, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return CreateCategoryOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
			},
		}
	} else if !user.Active {
		return CreateCategoryOutputDto{}, []util.ProblemDetails{
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

	newCategory, newCategoryErr := entities.NewCategory(user.ID, input.Name, input.Color)
	if newCategoryErr != nil {
		return CreateCategoryOutputDto{}, newCategoryErr
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
		CategoryID: newCategory.ID,
		Message:    "Category created successfully",
	}, nil
}
