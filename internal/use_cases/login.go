package usecases

import (
	"time"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/config"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
	"github.com/dgrijalva/jwt-go"
)

type LoginInputDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginOutputDto struct {
	Name           string `json:"name"`
	UserID         string `json:"user_id"`
	AccessToken    string `json:"access_token"`
	SuccessMessage string `json:"success_message"`
	ContentMessage string `json:"content_message"`
}

type LoginUseCase struct {
	UserRepository repositories.UserRepositoryInterface
}

func NewLoginUseCase(
	UserRepository repositories.UserRepositoryInterface,
) *LoginUseCase {
	return &LoginUseCase{
		UserRepository: UserRepository,
	}
}

func (c *LoginUseCase) Execute(input LoginInputDto) (LoginOutputDto, []util.ProblemDetails) {
	email, hashEmailWithHMACErr := util.HashEmailWithHMAC(input.Email)
	if hashEmailWithHMACErr != nil {
		return LoginOutputDto{}, hashEmailWithHMACErr
	}

	user, getUserByEmailErr := c.UserRepository.GetUserByEmail(email)
	if getUserByEmailErr != nil {
		return LoginOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Error getting user",
				Status:   500,
				Detail:   getUserByEmailErr.Error(),
				Instance: util.RFC500,
			},
		}
	} else if !user.Active {
		return LoginOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Forbidden",
				Title:    "User is not active",
				Status:   403,
				Detail:   "User is not active",
				Instance: util.RFC403,
			},
		}
	}

	if !user.Login.DecryptPassword(input.Password) {
		return LoginOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Unauthorized",
				Title:    "Invalid email or password",
				Status:   401,
				Detail:   "Invalid email or password",
				Instance: util.RFC401,
			},
		}
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := []byte(config.SECRETS_VAR.JWT_SECRET)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return LoginOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "JWT token Error",
				Status:   500,
				Detail:   "Error creating JWT token",
				Instance: util.RFC500,
			},
		}
	}

	return LoginOutputDto{
		Name:           user.Name,
		AccessToken:    tokenString,
		UserID:         user.ID,
		SuccessMessage: "Logged in successfully",
		ContentMessage: "Welcome, " + user.Name + "!",
	}, nil
}
