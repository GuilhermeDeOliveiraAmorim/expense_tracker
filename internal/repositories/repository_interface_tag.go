package repositories

import "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"

type TagRepositoryInterface interface {
	CreateTag(tag entities.Tag) error
	DeleteTag(tag entities.Tag) error
	GetCategories(userID string) ([]entities.Tag, error)
	GetTag(userID string, tagID string) (entities.Tag, error)
	ThisTagExists(userID string, tagName string) (bool, error)
}
