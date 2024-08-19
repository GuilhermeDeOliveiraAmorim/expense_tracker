package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
)

type GetUsersInputDto struct {
}

type GetUsersOutputDto struct {
	Users []entities.User `json:"users"`
}

type GetUsersUseCase struct {
	UserRepository repositories.UserRepositoryInterface
}

func NewGetUsersUseCase(
	UserRepository repositories.UserRepositoryInterface,
) *GetUsersUseCase {
	return &GetUsersUseCase{
		UserRepository: UserRepository,
	}
}

func (c *GetUsersUseCase) Execute(input GetUsersInputDto) (GetUsersOutputDto, []error) {
	searchedsUsers, err := c.UserRepository.GetUsers()
	if err != nil {
		return GetUsersOutputDto{}, err
	}

	return GetUsersOutputDto{
		Users: searchedsUsers,
	}, nil
}
