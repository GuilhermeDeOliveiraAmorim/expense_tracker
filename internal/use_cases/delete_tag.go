package usecases

import (
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type DeleteTagInputDto struct {
	UserID string `json:"user_id"`
	TagID  string `json:"tag_id"`
}

type DeleteTagOutputDto struct {
	Message string `json:"message"`
}

type DeleteTagUseCase struct {
	TagRepository  repositories.TagRepositoryInterface
	UserRepository repositories.UserRepositoryInterface
}

func NewDeleteTagUseCase(
	TagRepository repositories.TagRepositoryInterface,
	UserRepository repositories.UserRepositoryInterface,
) *DeleteTagUseCase {
	return &DeleteTagUseCase{
		TagRepository:  TagRepository,
		UserRepository: UserRepository,
	}
}

func (c *DeleteTagUseCase) Execute(input DeleteTagInputDto) (DeleteTagOutputDto, []util.ProblemDetails) {
	user, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return DeleteTagOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
			},
		}
	} else if !user.Active {
		return DeleteTagOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Forbidden",
				Title:    "User is not active",
				Status:   403,
				Detail:   "User is not active",
				Instance: util.RFC403,
			},
		}
	}

	tagToDelete, err := c.TagRepository.GetTag(input.UserID, input.TagID)
	if err != nil {
		return DeleteTagOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "Tag not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
			},
		}
	}

	tagToDelete.Deactivate()

	err = c.TagRepository.DeleteTag(tagToDelete)
	if err != nil {
		return DeleteTagOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Err deleting tag",
				Status:   500,
				Detail:   err.Error(),
				Instance: util.RFC500,
			},
		}
	}

	return DeleteTagOutputDto{
		Message: "Tag deleted successfully",
	}, nil
}
