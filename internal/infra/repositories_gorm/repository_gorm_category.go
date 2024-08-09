package repositoriesgorm

import (
	"errors"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	gorm *gorm.DB
}

func NewCategoryRepository(gorm *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		gorm: gorm,
	}
}

func (c *CategoryRepository) CreateCategory(category entities.Category) []error {
	if err := c.gorm.Create(&Categories{
		ID:            category.ID,
		Active:        category.Active,
		CreatedAt:     category.CreatedAt,
		UpdatedAt:     category.UpdatedAt,
		DeactivatedAt: category.DeactivatedAt,
		Name:          category.Name,
	}).Error; err != nil {
		return []error{err}
	}

	return nil
}

func (c *CategoryRepository) DeleteCategory(category entities.Category) []error {
	err := c.gorm.Model(&Categories{}).Where("id = ?", category.ID).Updates(Categories{
		Active:        category.Active,
		DeactivatedAt: category.DeactivatedAt,
	}).Error

	if err != nil {
		return []error{err}
	}

	return nil
}

func (c *CategoryRepository) GetCategories() ([]entities.Category, []error) {
	var categoriesModel []Categories
	if err := c.gorm.Find(&categoriesModel).Error; err != nil {
		return nil, []error{err}
	}

	var categories []entities.Category

	if len(categoriesModel) > 0 {
		for _, categoriesodel := range categoriesModel {
			category := entities.Category{
				SharedEntity: entities.SharedEntity{
					ID:            categoriesodel.ID,
					Active:        categoriesodel.Active,
					CreatedAt:     categoriesodel.CreatedAt,
					UpdatedAt:     categoriesodel.UpdatedAt,
					DeactivatedAt: categoriesodel.DeactivatedAt,
				},
				Name: categoriesodel.Name,
			}

			categories = append(categories, category)
		}
	}

	return categories, nil
}

func (c *CategoryRepository) GetCategory(categoryID string) (entities.Category, []error) {
	var categoryModel Categories

	result := c.gorm.Model(&Categories{}).Where("id = ?", categoryID).First(&categoryModel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return entities.Category{}, []error{errors.New("category not found")}
		}
		return entities.Category{}, []error{errors.New(result.Error.Error())}
	}

	category := entities.Category{
		SharedEntity: entities.SharedEntity{
			ID:            categoryModel.ID,
			Active:        categoryModel.Active,
			CreatedAt:     categoryModel.CreatedAt,
			UpdatedAt:     categoryModel.UpdatedAt,
			DeactivatedAt: categoryModel.DeactivatedAt,
		},
		Name: categoryModel.Name,
	}

	return category, nil
}

func (c *CategoryRepository) UpdateCategory(category entities.Category) []error {
	result := c.gorm.Model(&Categories{}).Where("id", category.ID).Updates(Categories{
		Name:      category.Name,
		UpdatedAt: category.UpdatedAt,
	})

	if result.Error != nil {
		return []error{errors.New(result.Error.Error())}
	}

	return nil
}
