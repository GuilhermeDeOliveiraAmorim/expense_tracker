package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
)

type CreateUserInputDto struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserOutputDto struct {
	ID string `json:"id"`
}

type CreateUserUseCase struct {
	UserRepository repositories.UserRepositoryInterface
}

func (c *CreateUserUseCase) Execute(input CreateUserInputDto) (CreateUserOutputDto, []error) {
	newLogin, err := entities.NewLogin(input.Email, input.Password)
	if err != nil {
		return CreateUserOutputDto{}, err
	}

	newLogin.EncryptEmail()
	newLogin.EncryptPassword()

	newUser, err := entities.NewUser(input.Name, *newLogin)
	if err != nil {
		return CreateUserOutputDto{}, err
	}

	errs := c.UserRepository.CreateUser(*newUser)
	if errs != nil {
		return CreateUserOutputDto{}, err
	}

	return CreateUserOutputDto{
		ID: newUser.ID,
	}, nil
}
