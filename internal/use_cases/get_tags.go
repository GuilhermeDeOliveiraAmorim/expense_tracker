package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type GetTagsInputDto struct {
	UserID string `json:"user_id"`
}

type GetTagsOutputDto struct {
	Tags []entities.Tag `json:"tags"`
}

type GetTagsUseCase struct {
	TagRepository repositories.TagRepositoryInterface
	UserRepository     repositories.UserRepositoryInterface
}

func NewGetTagsUseCase(
	TagRepository repositories.TagRepositoryInterface,
	UserRepository repositories.UserRepositoryInterface,
) *GetTagsUseCase {
	return &GetTagsUseCase{
		TagRepository: TagRepository,
		UserRepository:     UserRepository,
	}
}

func (c *GetTagsUseCase) Execute(input GetTagsInputDto) (GetTagsOutputDto, []util.ProblemDetails) {
	user, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return GetTagsOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
			},
		}
	} else if !user.Active {
		return GetTagsOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Forbidden",
				Title:    "User is not active",
				Status:   403,
				Detail:   "User is not active",
				Instance: util.RFC403,
			},
		}
	}

	searchedsTags, err := c.TagRepository.GetTags(input.UserID)
	if err != nil {
		return GetTagsOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Err fetching tags",
				Status:   500,
				Detail:   err.Error(),
				Instance: util.RFC500,
			},
		}
	}

	return GetTagsOutputDto{
		Tags: searchedsTags,
	}, nil
}
