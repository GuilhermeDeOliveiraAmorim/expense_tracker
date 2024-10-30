package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type DeleteUserInputDto struct {
	UserID string `json:"user_id"`
}

type DeleteUserOutputDto struct {
	SuccessMessage string `json:"success_message"`
	ContentMessage string `json:"content_message"`
}

type DeleteUserUseCase struct {
	UserRepository repositories.UserRepositoryInterface
}

func NewDeleteUserUseCase(
	UserRepository repositories.UserRepositoryInterface,
) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		UserRepository: UserRepository,
	}
}

func (c *DeleteUserUseCase) Execute(input DeleteUserInputDto) (DeleteUserOutputDto, []util.ProblemDetails) {
	userToDelete, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return DeleteUserOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   "User not found",
				Instance: util.RFC404,
			},
		}
	}

	userToDelete.Deactivate()

	DeleteUserErr := c.UserRepository.DeleteUser(userToDelete)
	if DeleteUserErr != nil {
		return DeleteUserOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Err deleting user",
				Status:   500,
				Detail:   DeleteUserErr.Error(),
				Instance: util.RFC500,
			},
		}
	}

	return DeleteUserOutputDto{
		SuccessMessage: "User deleted successfully",
		ContentMessage: "User " + userToDelete.Name + " deleted",
	}, nil
}
