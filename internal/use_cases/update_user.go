package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type UpdateUserInputDto struct {
	UserID string `json:"id"`
	Name   string `json:"name"`
}

type UpdateUserOutputDto struct {
	ID string `json:"id"`
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
		ID: searchedUser.ID,
	}, nil
}
