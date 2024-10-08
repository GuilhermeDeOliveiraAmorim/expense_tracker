package factory

import (
	repositoriesgorm "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/infra/repositories_gorm"
	usecases "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/use_cases"
	"gorm.io/gorm"
)

type CategoryFactory struct {
	CreateCategory *usecases.CreateCategoryUseCase
	DeleteCategory *usecases.DeleteCategoryUseCase
	GetCategories  *usecases.GetCategoriesUseCase
	GetCategory    *usecases.GetCategoryUseCase
	UpdateCategory *usecases.UpdateCategoryUseCase
}

func NewCategoryFactory(db *gorm.DB) *CategoryFactory {
	categoryRepository := repositoriesgorm.NewCategoryRepository(db)
	userRepository := repositoriesgorm.NewUserRepository(db)

	createCategory := usecases.NewCreateCategoryUseCase(categoryRepository, userRepository)
	deleteCategory := usecases.NewDeleteCategoryUseCase(categoryRepository, userRepository)
	getCategories := usecases.NewGetCategoriesUseCase(categoryRepository, userRepository)
	getCategory := usecases.NewGetCategoryUseCase(categoryRepository, userRepository)
	updateCategory := usecases.NewUpdateCategoryUseCase(categoryRepository, userRepository)

	return &CategoryFactory{
		CreateCategory: createCategory,
		DeleteCategory: deleteCategory,
		GetCategories:  getCategories,
		GetCategory:    getCategory,
		UpdateCategory: updateCategory,
	}
}
