package repositories

import "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"

type CategoryRepositoryInterface interface {
	CreateCategory(category entities.Category) error
	DeleteCategory(category entities.Category) error
	GetCategories() ([]entities.Category, error)
	GetCategory(categoryID string) (entities.Category, error)
	ThisCategoryExists(categoryName string) (bool, error)
	UpdateCategory(category entities.Category) error
}
