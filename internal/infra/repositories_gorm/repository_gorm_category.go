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

func (c *CategoryRepository) Create(category *entities.Category) error {
	if err := c.gorm.Create(&Categories{
		ID:            category.ID,
		Active:        category.Active,
		CreatedAt:     category.CreatedAt,
		UpdatedAt:     category.UpdatedAt,
		DeactivatedAt: category.DeactivatedAt,
		Name:          category.Name,
	}).Error; err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func (c *CategoryRepository) Delete(category *entities.Category) error {
	err := c.gorm.Model(&Categories{}).Where("id = ?", category.ID).Updates(Categories{
		Active:        category.Active,
		DeactivatedAt: category.DeactivatedAt,
	}).Error

	if err != nil {
		return err
	}

	return nil
}

func (c *CategoryRepository) Update(category *entities.Category) error {
	result := c.gorm.Model(&Categories{}).Where("id", category.ID).Updates(Categories{
		Name:      category.Name,
		UpdatedAt: category.UpdatedAt,
	})

	if result.Error != nil {
		return errors.New(result.Error.Error())
	}

	return nil
}

func (c *CategoryRepository) GetByID(categoryID string) (entities.Category, error) {
	var categoryModel Categories

	result := c.gorm.Model(&Categories{}).Where("id = ?", categoryID).First(&categoryModel)
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
		Name: categoryModel.Name,
	}

	return category, nil
}
