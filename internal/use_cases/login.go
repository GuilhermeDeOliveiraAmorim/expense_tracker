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

	user, err := c.UserRepository.GetUserByEmail(email)
	if err != nil {
		return LoginOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Error getting user",
				Status:   500,
				Detail:   "Error getting user by email",
				Instance: util.RFC500,
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
		"email":   user.Login.Email,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := []byte(config.SECRETS_LOCAL.JWT_SECRET)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return LoginOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Error creating JWT token",
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
