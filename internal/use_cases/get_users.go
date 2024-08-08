package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
)

type GetUserInputDto struct {
	UserID string `json:"user_id"`
}

type GetUserOutputDto struct {
	User entities.User `json:"user"`
}

type GetUserUseCase struct {
	UserRepository repositories.UserRepositoryInterface
}

func (c *GetUserUseCase) Execute(input GetUserInputDto) (GetUserOutputDto, []error) {
	searchedUser, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return GetUserOutputDto{}, err
	}

	return GetUserOutputDto{
		User: searchedUser,
	}, nil
}
