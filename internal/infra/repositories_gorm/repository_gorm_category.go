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

func (cr *CategoryRepository) Delete(category *entities.Category) error {
	err := cr.gorm.Model(&Categories{}).Where("id = ?", category.ID).Updates(map[string]interface{}{
		"Active":        category.Active,
		"DeactivatedAt": category.DeactivatedAt,
	}).Error

	if err != nil {
		return err
	}

	return nil
}
