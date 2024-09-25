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

func (c *CategoryRepository) CreateCategory(category entities.Category) error {
	if err := c.gorm.Create(&Categories{
		ID:            category.ID,
		Active:        category.Active,
		CreatedAt:     category.CreatedAt,
		UpdatedAt:     category.UpdatedAt,
		DeactivatedAt: category.DeactivatedAt,
		UserID:        category.UserID,
		Name:          category.Name,
		Color:         category.Color,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (c *CategoryRepository) DeleteCategory(category entities.Category) error {
	result := c.gorm.Model(&Categories{}).Where("id = ? AND user_id = ? AND active = ?", category.ID, category.UserID, true).
		Select("Active", "DeactivatedAt", "UpdatedAt").Updates(Categories{
		Active:        category.Active,
		DeactivatedAt: category.DeactivatedAt,
		UpdatedAt:     category.UpdatedAt,
	})

	if result.Error != nil {
		return errors.New(result.Error.Error())
	}

	return nil
}

func (c *CategoryRepository) GetCategories(userID string) ([]entities.Category, error) {
	var categoriesModel []Categories
	if err := c.gorm.Where("user_id = ? AND active = ?", userID, true).Find(&categoriesModel).Error; err != nil {
		return nil, err
	}

	var categories []entities.Category

	if len(categoriesModel) > 0 {
		for _, categoryModel := range categoriesModel {
			category := entities.Category{
				SharedEntity: entities.SharedEntity{
					ID:            categoryModel.ID,
					Active:        categoryModel.Active,
					CreatedAt:     categoryModel.CreatedAt,
					UpdatedAt:     categoryModel.UpdatedAt,
					DeactivatedAt: categoryModel.DeactivatedAt,
				},
				UserID: categoryModel.UserID,
				Name:   categoryModel.Name,
				Color:  categoryModel.Color,
			}

			categories = append(categories, category)
		}
	}

	return categories, nil
}

func (c *CategoryRepository) GetCategory(userID string, categoryID string) (entities.Category, error) {
	var categoryModel Categories

	result := c.gorm.Model(&Categories{}).Where("id = ? AND user_id = ? AND active = ?", categoryID, userID, true).First(&categoryModel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return entities.Category{}, errors.New("category not found")
		}
		return entities.Category{}, errors.New(result.Error.Error())
	}

	category := entities.Category{
		SharedEntity: entities.SharedEntity{
			ID:            categoryModel.ID,
			Active:        categoryModel.Active,
			CreatedAt:     categoryModel.CreatedAt,
			UpdatedAt:     categoryModel.UpdatedAt,
			DeactivatedAt: categoryModel.DeactivatedAt,
		},
		UserID: categoryModel.UserID,
		Name:   categoryModel.Name,
		Color:  categoryModel.Color,
	}

	return category, nil
}

func (c *CategoryRepository) UpdateCategory(category entities.Category) error {
	result := c.gorm.Model(&Categories{}).Where("id = ? AND user_id = ? AND active = ?", category.ID, category.UserID, true).Updates(Categories{
		Name:      category.Name,
		Color:     category.Color,
		UpdatedAt: category.UpdatedAt,
	})

	if result.Error != nil {
		return errors.New(result.Error.Error())
	}

	return nil
}

func (c *CategoryRepository) ThisCategoryExists(userID string, categoryName string) (bool, error) {
	var categoryModel Categories

	result := c.gorm.Model(&Categories{}).Where("name = ? AND user_id = ? AND active = ?", categoryName, userID, true).First(&categoryModel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, errors.New("category not found")
		}
		return false, errors.New(result.Error.Error())
	}

	return true, nil
}
