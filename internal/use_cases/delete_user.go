package usecases

import "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"

type DeleteUserInputDto struct {
	UserID string `json:"id"`
}

type DeleteUserOutputDto struct {
	ID string `json:"id"`
}

type DeleteUserUseCase struct {
	UserRepository repositories.UserRepositoryInterface
}

func (c *DeleteUserUseCase) Execute(input DeleteUserInputDto) (DeleteUserOutputDto, []error) {
	userToDelete, errs := c.UserRepository.GetUser(input.UserID)
	if errs != nil {
		return DeleteUserOutputDto{}, errs
	}

	userToDelete.Deactivate()

	errs = c.UserRepository.DeleteUser(userToDelete)
	if errs != nil {
		return DeleteUserOutputDto{}, errs
	}

	return DeleteUserOutputDto{
		ID: userToDelete.ID,
	}, nil
}
