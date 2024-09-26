package usecases

import (
	"strings"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type UpdateTagInputDto struct {
	UserID string `json:"user_id"`
	TagID  string `json:"tag_id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
}

type UpdateTagOutputDto struct {
	TagID          string `json:"tag_id"`
	SuccessMessage string `json:"success_message"`
	ContentMessage string `json:"content_message"`
}

type UpdateTagUseCase struct {
	TagRepository  repositories.TagRepositoryInterface
	UserRepository repositories.UserRepositoryInterface
}

func NewUpdateTagUseCase(
	TagRepository repositories.TagRepositoryInterface,
	UserRepository repositories.UserRepositoryInterface,
) *UpdateTagUseCase {
	return &UpdateTagUseCase{
		TagRepository:  TagRepository,
		UserRepository: UserRepository,
	}
}

func (c *UpdateTagUseCase) Execute(input UpdateTagInputDto) (UpdateTagOutputDto, []util.ProblemDetails) {
	user, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return UpdateTagOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
			},
		}
	} else if !user.Active {
		return UpdateTagOutputDto{}, []util.ProblemDetails{
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
		return UpdateTagOutputDto{}, []util.ProblemDetails{
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
		return UpdateTagOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Validation Error",
				Title:    "Tag already exists",
				Status:   409,
				Detail:   "A tag with this name already exists",
				Instance: util.RFC409,
			},
		}
	}

	searchedTag, getTagErr := c.TagRepository.GetTag(input.UserID, input.TagID)
	if getTagErr != nil {
		return UpdateTagOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "Tag not found",
				Status:   404,
				Detail:   getTagErr.Error(),
				Instance: util.RFC400,
			},
		}
	}

	changeNameErr := searchedTag.ChangeName(input.Name)
	if len(changeNameErr) > 0 {
		return UpdateTagOutputDto{}, changeNameErr
	}

	changeColorErr := searchedTag.ChangeColor(input.Color)
	if len(changeColorErr) > 0 {
		return UpdateTagOutputDto{}, changeColorErr
	}

	updateTagErr := c.TagRepository.UpdateTag(searchedTag)
	if updateTagErr != nil {
		return UpdateTagOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Validation Error",
				Title:    "Bad Request",
				Status:   500,
				Detail:   updateTagErr.Error(),
				Instance: util.RFC500,
			},
		}
	}

	return UpdateTagOutputDto{
		TagID:          searchedTag.ID,
		SuccessMessage: "Tag updated successfully",
		ContentMessage: "Tag ID: " + searchedTag.ID,
	}, nil
}
