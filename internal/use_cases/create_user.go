package usecases

import (
	"strings"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
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

func NewCreateUserUseCase(
	UserRepository repositories.UserRepositoryInterface,
) *CreateUserUseCase {
	return &CreateUserUseCase{
		UserRepository: UserRepository,
	}
}

func (c *CreateUserUseCase) Execute(input CreateUserInputDto) (CreateUserOutputDto, []util.ProblemDetails) {
	existingUser, GetUserByNameErr := c.UserRepository.ThisUserExists(input.Name)
	if GetUserByNameErr != nil && strings.Compare(GetUserByNameErr.Error(), "user not found") > 0 {
		return CreateUserOutputDto{}, []util.ProblemDetails{
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
		return CreateUserOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Validation Error",
				Title:    "User already exists",
				Status:   409,
				Detail:   "A user with this name already exists",
				Instance: util.RFC409,
			},
		}
	}

	newLogin, err := entities.NewLogin(input.Email, input.Password)
	if err != nil {
		return CreateUserOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Validation Error",
				Title:    "Invalid email or password",
				Status:   400,
				Detail:   "Email or password is invalid",
				Instance: util.RFC400,
			},
		}
	}

	EncryptEmailErr := newLogin.EncryptEmail()
	if EncryptEmailErr != nil {
		return CreateUserOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Error encrypting email",
				Status:   500,
				Detail:   EncryptEmailErr.Error(),
				Instance: util.RFC500,
			},
		}
	}

	EncryptPasswordErr := newLogin.EncryptPassword()
	if EncryptPasswordErr != nil {
		return CreateUserOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Error encrypting password",
				Status:   500,
				Detail:   EncryptPasswordErr.Error(),
				Instance: util.RFC500,
			},
		}
	}

	newUser, err := entities.NewUser(input.Name, *newLogin)
	if err != nil {
		return CreateUserOutputDto{}, err
	}

	CreateUserErr := c.UserRepository.CreateUser(*newUser)
	if CreateUserErr != nil {
		return CreateUserOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Error creating new user",
				Status:   500,
				Detail:   CreateUserErr.Error(),
				Instance: util.RFC500,
			},
		}
	}

	return CreateUserOutputDto{
		ID: newUser.ID,
	}, nil
}
