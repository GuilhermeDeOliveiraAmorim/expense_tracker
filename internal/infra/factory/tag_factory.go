package factory

import (
	repositoriesgorm "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/infra/repositories_gorm"
	usecases "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/use_cases"
	"gorm.io/gorm"
)

type TagFactory struct {
	CreateTag *usecases.CreateTagUseCase
	DeleteTag *usecases.DeleteTagUseCase
	GetTags   *usecases.GetTagsUseCase
	GetTag    *usecases.GetTagUseCase
}

func NewTagFactory(db *gorm.DB) *TagFactory {
	tagRepository := repositoriesgorm.NewTagRepository(db)
	userRepository := repositoriesgorm.NewUserRepository(db)

	createTag := usecases.NewCreateTagUseCase(tagRepository, userRepository)
	deleteTag := usecases.NewDeleteTagUseCase(tagRepository, userRepository)
	getTags := usecases.NewGetTagsUseCase(tagRepository, userRepository)
	getTag := usecases.NewGetTagUseCase(tagRepository, userRepository)

	return &TagFactory{
		CreateTag: createTag,
		DeleteTag: deleteTag,
		GetTags:   getTags,
		GetTag:    getTag,
	}
}
