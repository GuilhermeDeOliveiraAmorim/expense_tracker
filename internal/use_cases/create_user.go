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
	Name           string `json:"name"`
	UserID         string `json:"user_id"`
	SuccessMessage string `json:"success_message"`
	ContentMessage string `json:"content_message"`
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
	email, hashEmailWithHMACErr := util.HashEmailWithHMAC(input.Email)
	if hashEmailWithHMACErr != nil {
		return CreateUserOutputDto{}, hashEmailWithHMACErr
	}

	userEmailExists, userEmailExistsErr := c.UserRepository.ThisUserEmailExists(email)
	if userEmailExists {
		return CreateUserOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Conflict",
				Title:    "Email already exists",
				Status:   409,
				Detail:   "Email already exists",
				Instance: util.RFC409,
			},
		}
	} else if strings.Compare(userEmailExistsErr.Error(), "not found") != 0 {
		return CreateUserOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Error checking user email existence",
				Status:   500,
				Detail:   userEmailExistsErr.Error(),
				Instance: util.RFC500,
			},
		}
	}

	userNameExists, userNameExistsErr := c.UserRepository.ThisUserNameExists(input.Name)
	if userNameExists {
		return CreateUserOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Conflict",
				Title:    "Username already exists",
				Status:   409,
				Detail:   "Username already exists",
				Instance: util.RFC409,
			},
		}
	} else if strings.Compare(userNameExistsErr.Error(), "not found") != 0 {
		return CreateUserOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Error checking user name existence",
				Status:   500,
				Detail:   userEmailExistsErr.Error(),
				Instance: util.RFC500,
			},
		}
	}

	newLogin, newLoginErr := entities.NewLogin(input.Email, input.Password)
	if newLoginErr != nil {
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

	encryptEmailErr := newLogin.EncryptEmail()
	if encryptEmailErr != nil {
		return CreateUserOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Error encrypting email",
				Status:   500,
				Detail:   encryptEmailErr.Error(),
				Instance: util.RFC500,
			},
		}
	}

	encryptPasswordErr := newLogin.EncryptPassword()
	if encryptPasswordErr != nil {
		return CreateUserOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Error encrypting password",
				Status:   500,
				Detail:   encryptPasswordErr.Error(),
				Instance: util.RFC500,
			},
		}
	}

	newUser, newUserErr := entities.NewUser(input.Name, *newLogin)
	if newUserErr != nil {
		return CreateUserOutputDto{}, newUserErr
	}

	createUserErr := c.UserRepository.CreateUser(*newUser)
	if createUserErr != nil {
		return CreateUserOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Error creating new user",
				Status:   500,
				Detail:   createUserErr.Error(),
				Instance: util.RFC500,
			},
		}
	}

	return CreateUserOutputDto{
		Name:           newUser.Name,
		UserID:         newUser.ID,
		SuccessMessage: "User created successfully",
		ContentMessage: "Welcome, " + newUser.Name + "!",
	}, nil
}
