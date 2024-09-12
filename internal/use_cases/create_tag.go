package usecases

import (
	"strings"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type CreateTagInputDto struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
}

type CreateTagOutputDto struct {
	TagID   string `json:"tag_id"`
	Message string `json:"message"`
}

type CreateTagUseCase struct {
	TagRepository  repositories.TagRepositoryInterface
	UserRepository repositories.UserRepositoryInterface
}

func NewCreateTagUseCase(
	TagRepository repositories.TagRepositoryInterface,
	UserRepository repositories.UserRepositoryInterface,
) *CreateTagUseCase {
	return &CreateTagUseCase{
		TagRepository:  TagRepository,
		UserRepository: UserRepository,
	}
}

func (c *CreateTagUseCase) Execute(input CreateTagInputDto) (CreateTagOutputDto, []util.ProblemDetails) {
	user, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return CreateTagOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
			},
		}
	} else if !user.Active {
		return CreateTagOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Forbidden",
				Title:    "User is not active",
				Status:   403,
				Detail:   "User is not active",
				Instance: util.RFC403,
			},
		}
	}

	existingTag, GetTagByNameErr := c.TagRepository.ThisTagExists(input.UserID, input.Name)
	if GetTagByNameErr != nil && strings.Compare(GetTagByNameErr.Error(), "tag not found") > 0 {
		return CreateTagOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Error fetching existing tag",
				Status:   500,
				Detail:   GetTagByNameErr.Error(),
				Instance: util.RFC500,
			},
		}
	}

	if existingTag {
		return CreateTagOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Validation Error",
				Title:    "Tag already exists",
				Status:   409,
				Detail:   "A tag with this name already exists",
				Instance: util.RFC409,
			},
		}
	}

	newTag, newTagErr := entities.NewTag(user.ID, input.Name, input.Color)
	if newTagErr != nil {
		return CreateTagOutputDto{}, newTagErr
	}

	CreateTagErr := c.TagRepository.CreateTag(*newTag)
	if CreateTagErr != nil {
		return CreateTagOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Error creating new tag",
				Status:   500,
				Detail:   CreateTagErr.Error(),
				Instance: util.RFC500,
			},
		}
	}

	return CreateTagOutputDto{
		TagID:   newTag.ID,
		Message: "Tag created successfully",
	}, nil
}
