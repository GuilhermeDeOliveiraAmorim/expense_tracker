package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type GetCategoriesInputDto struct {
	UserID string `json:"user_id"`
}

type GetCategoriesOutputDto struct {
	Categories []entities.Category `json:"categories"`
}

type GetCategoriesUseCase struct {
	CategoryRepository repositories.CategoryRepositoryInterface
	UserRepository     repositories.UserRepositoryInterface
}

func NewGetCategoriesUseCase(
	CategoryRepository repositories.CategoryRepositoryInterface,
	UserRepository repositories.UserRepositoryInterface,
) *GetCategoriesUseCase {
	return &GetCategoriesUseCase{
		CategoryRepository: CategoryRepository,
		UserRepository:     UserRepository,
	}
}

func (c *GetCategoriesUseCase) Execute(input GetCategoriesInputDto) (GetCategoriesOutputDto, []util.ProblemDetails) {
	user, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return GetCategoriesOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
			},
		}
	} else if !user.Active {
		return GetCategoriesOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Forbidden",
				Title:    "User is not active",
				Status:   403,
				Detail:   "User is not active",
				Instance: util.RFC403,
			},
		}
	}

	searchedsCategories, err := c.CategoryRepository.GetCategories(input.UserID)
	if err != nil {
		return GetCategoriesOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Err fetching categories",
				Status:   500,
				Detail:   err.Error(),
				Instance: util.RFC500,
			},
		}
	}

	return GetCategoriesOutputDto{
		Categories: searchedsCategories,
	}, nil
}
