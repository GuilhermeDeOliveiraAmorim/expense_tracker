package usecases

import (
	"strings"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type UpdateUserInputDto struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

type UpdateUserOutputDto struct {
	UserID         string `json:"user_id"`
	SuccessMessage string `json:"success_message"`
	ContentMessage string `json:"content_message"`
}

type UpdateUserUseCase struct {
	UserRepository repositories.UserRepositoryInterface
}

func NewUpdateUserUseCase(
	UserRepository repositories.UserRepositoryInterface,
) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		UserRepository: UserRepository,
	}
}

func (c *UpdateUserUseCase) Execute(input UpdateUserInputDto) (UpdateUserOutputDto, []util.ProblemDetails) {
	user, getUserErr := c.UserRepository.GetUser(input.UserID)
	if getUserErr != nil {
		return UpdateUserOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   getUserErr.Error(),
				Instance: util.RFC404,
			},
		}
	} else if !user.Active {
		return UpdateUserOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Forbidden",
				Title:    "User is not active",
				Status:   403,
				Detail:   "User is not active",
				Instance: util.RFC403,
			},
		}
	}

	existingUser, GetUserByNameErr := c.UserRepository.ThisUserExists(input.Name)
	if GetUserByNameErr != nil && strings.Compare(GetUserByNameErr.Error(), "user not found") > 0 {
		return UpdateUserOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Error fetching existing user",
				Status:   500,
				Detail:   GetUserByNameErr.Error(),
				Instance: util.RFC500,
			},
		}
	}

	if existingUser {
		return UpdateUserOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Validation Error",
				Title:    "User already exists",
				Status:   409,
				Detail:   "A user with this name already exists",
				Instance: util.RFC409,
			},
		}
	}

	searchedUser, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return UpdateUserOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
			},
		}
	}

	if input.Name != searchedUser.Name {
		err := searchedUser.ChangeName(input.Name)
		if len(err) > 0 {
			return UpdateUserOutputDto{}, err
		}

		UpdateUserErr := c.UserRepository.UpdateUser(searchedUser)
		if UpdateUserErr != nil {
			return UpdateUserOutputDto{}, err
		}
	} else {
		return UpdateUserOutputDto{}, []util.ProblemDetails{
			{
				Type:     "No Changes Made",
				Title:    "No changes detected",
				Status:   204,
				Detail:   "No changes were made to the user",
				Instance: util.RFC204,
			},
		}
	}

	return UpdateUserOutputDto{
		UserID:         searchedUser.ID,
		SuccessMessage: "User updated successfully",
		ContentMessage: "Your new name is " + searchedUser.Name + "!",
	}, nil
}
