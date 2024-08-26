package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type UserOutput struct {
	entities.SharedEntity
	Name string `json:"name"`
}

type GetUserInputDto struct {
	UserID string `json:"user_id"`
}

type GetUserOutputDto struct {
	User UserOutput `json:"user"`
}

type GetUserUseCase struct {
	UserRepository repositories.UserRepositoryInterface
}

func NewGetUserUseCase(
	UserRepository repositories.UserRepositoryInterface,
) *GetUserUseCase {
	return &GetUserUseCase{
		UserRepository: UserRepository,
	}
}

func (c *GetUserUseCase) Execute(input GetUserInputDto) (GetUserOutputDto, []util.ProblemDetails) {
	searchedUser, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return GetUserOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   "User not found",
				Instance: util.RFC404,
			},
		}
	}

	return GetUserOutputDto{
		User: UserOutput{
			SharedEntity: searchedUser.SharedEntity,
			Name:         searchedUser.Name,
		},
	}, nil
}
