package repositoriesgorm

import (
	"errors"
	"sort"

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
	tx := c.gorm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := tx.Create(&Categories{
		ID:            category.ID,
		Active:        category.Active,
		CreatedAt:     category.CreatedAt,
		UpdatedAt:     category.UpdatedAt,
		DeactivatedAt: category.DeactivatedAt,
		UserID:        category.UserID,
		Name:          category.Name,
		Color:         category.Color,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (c *CategoryRepository) DeleteCategory(category entities.Category) error {
	tx := c.gorm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var expenseCount int64
	if err := tx.Model(&Expenses{}).Where("category_id = ? AND active = ?", category.ID, true).Count(&expenseCount).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to check expenses associated with the category: " + err.Error())
	}

	if expenseCount > 0 {
		tx.Rollback()
		return errors.New("there are expenses associated with this category")
	}

	result := tx.Model(&Categories{}).Where("id = ? AND user_id = ? AND active = ?", category.ID, category.UserID, true).
		Select("Active", "DeactivatedAt", "UpdatedAt").Updates(Categories{
		Active:        category.Active,
		DeactivatedAt: category.DeactivatedAt,
		UpdatedAt:     category.UpdatedAt,
	})

	if result.Error != nil {
		tx.Rollback()
		return errors.New(result.Error.Error())
	}

	return tx.Commit().Error
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

		sort.Slice(categories, func(i, j int) bool {
			return categories[i].CreatedAt.After(categories[j].CreatedAt)
		})
	} else {
		categories = []entities.Category{}
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
	tx := c.gorm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	result := tx.Model(&Categories{}).Where("id = ? AND user_id = ? AND active = ?", category.ID, category.UserID, true).Updates(Categories{
		Name:      category.Name,
		Color:     category.Color,
		UpdatedAt: category.UpdatedAt,
	})

	if result.Error != nil {
		tx.Rollback()
		return errors.New(result.Error.Error())
	}

	return tx.Commit().Error
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
