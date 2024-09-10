package repositories

import "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"

type CategoryRepositoryInterface interface {
	CreateCategory(category entities.Category) error
	DeleteCategory(category entities.Category) error
	GetCategories(userID string) ([]entities.Category, error)
	GetCategory(userID string, categoryID string) (entities.Category, error)
	ThisCategoryExists(userID string, categoryName string) (bool, error)
	UpdateCategory(category entities.Category) error
}
