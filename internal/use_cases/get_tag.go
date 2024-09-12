package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type GetTagInputDto struct {
	UserID string `json:"user_id"`
	TagID  string `json:"tag_id"`
}

type GetTagOutputDto struct {
	Tag entities.Tag `json:"tag"`
}

type GetTagUseCase struct {
	TagRepository  repositories.TagRepositoryInterface
	UserRepository repositories.UserRepositoryInterface
}

func NewGetTagUseCase(
	TagRepository repositories.TagRepositoryInterface,
	UserRepository repositories.UserRepositoryInterface,
) *GetTagUseCase {
	return &GetTagUseCase{
		TagRepository:  TagRepository,
		UserRepository: UserRepository,
	}
}

func (c *GetTagUseCase) Execute(input GetTagInputDto) (GetTagOutputDto, []util.ProblemDetails) {
	user, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return GetTagOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
			},
		}
	} else if !user.Active {
		return GetTagOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Forbidden",
				Title:    "User is not active",
				Status:   403,
				Detail:   "User is not active",
				Instance: util.RFC403,
			},
		}
	}

	searchedTag, err := c.TagRepository.GetTag(input.UserID, input.TagID)
	if err != nil {
		return GetTagOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "Tag not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC500,
			},
		}
	}

	return GetTagOutputDto{
		Tag: searchedTag,
	}, nil
}
