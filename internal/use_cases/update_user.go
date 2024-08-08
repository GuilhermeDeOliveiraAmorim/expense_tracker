package usecases

import (
	"fmt"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
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

func (c *UpdateUserUseCase) Execute(input UpdateUserInputDto) (UpdateUserOutputDto, []error) {
	searchedUser, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return UpdateUserOutputDto{}, err
	}

	if input.Name != searchedUser.Name {
		err = searchedUser.ChangeName(input.Name)
		if len(err) > 0 {
			return UpdateUserOutputDto{}, err
		}

		err = c.UserRepository.UpdateUser(searchedUser)
		if err != nil {
			return UpdateUserOutputDto{}, err
		}
	} else {
		return UpdateUserOutputDto{}, []error{
			fmt.Errorf("name cannot be the same as the current one"),
		}
	}

	return UpdateUserOutputDto{
		ID: searchedUser.ID,
	}, nil
}
