package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type GetCategoryInputDto struct {
	UserID     string `json:"user_id"`
	CategoryID string `json:"category_id"`
}

type GetCategoryOutputDto struct {
	Category entities.Category `json:"category"`
}

type GetCategoryUseCase struct {
	CategoryRepository repositories.CategoryRepositoryInterface
}

func NewGetCategoryUseCase(
	CategoryRepository repositories.CategoryRepositoryInterface,
) *GetCategoryUseCase {
	return &GetCategoryUseCase{
		CategoryRepository: CategoryRepository,
	}
}

func (c *GetCategoryUseCase) Execute(input GetCategoryInputDto) (GetCategoryOutputDto, []util.ProblemDetails) {
	searchedCategory, err := c.CategoryRepository.GetCategory(input.UserID, input.CategoryID)
	if err != nil {
		return GetCategoryOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "Category not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC500,
			},
		}
	}

	return GetCategoryOutputDto{
		Category: searchedCategory,
	}, nil
}
