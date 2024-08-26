package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
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

func (c *GetUsersUseCase) Execute(input GetUsersInputDto) (GetUsersOutputDto, []util.ProblemDetails) {
	searchedsUsers, err := c.UserRepository.GetUsers()
	if err != nil {
		return GetUsersOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Error fetching users",
				Status:   500,
				Detail:   err.Error(),
				Instance: util.RFC500,
			},
		}
	}

	return GetUsersOutputDto{
		Users: searchedsUsers,
	}, nil
}
