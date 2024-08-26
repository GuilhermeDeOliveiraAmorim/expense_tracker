package repositoriesgorm

import (
	"errors"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"gorm.io/gorm"
)

type ExpenseRepository struct {
	gorm *gorm.DB
}

func NewExpenseRepository(gorm *gorm.DB) *ExpenseRepository {
	return &ExpenseRepository{
		gorm: gorm,
	}
}

func (e *ExpenseRepository) CreateExpense(expense entities.Expense) error {
	if err := e.gorm.Create(&Expenses{
		ID:            expense.ID,
		Active:        expense.Active,
		CreatedAt:     expense.CreatedAt,
		UpdatedAt:     expense.UpdatedAt,
		DeactivatedAt: expense.DeactivatedAt,
		UserID:        expense.UserID,
		Amount:        expense.Amount,
		CategoryID:    expense.Category.ID,
		Notes:         expense.Notes,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (e *ExpenseRepository) DeleteExpense(category entities.Expense) error {
	err := e.gorm.Model(&Expenses{}).Where("id = ?", category.ID).Updates(Expenses{
		Active:        category.Active,
		DeactivatedAt: category.DeactivatedAt,
	}).Error

	if err != nil {
		return err
	}

	return nil
}

func (e *ExpenseRepository) GetExpenses() ([]entities.Expense, error) {
	var expensesModel []Expenses

	if err := e.gorm.Find(&expensesModel).Error; err != nil {
		return nil, err
	}

	var expenses []entities.Expense

	if len(expensesModel) > 0 {
		for _, expenseModel := range expensesModel {
			var categoryModel Categories

			result := e.gorm.Model(&Categories{}).Where("id = ?", expenseModel.CategoryID).First(&categoryModel)
			if result.Error != nil {
				if errors.Is(result.Error, gorm.ErrRecordNotFound) {
					return []entities.Expense{}, errors.New("error searching expense (" + expenseModel.ID + ") category")
				}
				return []entities.Expense{}, errors.New(result.Error.Error())
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

			expense := entities.Expense{
				SharedEntity: entities.SharedEntity{
					ID:            expenseModel.ID,
					Active:        expenseModel.Active,
					CreatedAt:     expenseModel.CreatedAt,
					UpdatedAt:     expenseModel.UpdatedAt,
					DeactivatedAt: expenseModel.DeactivatedAt,
				},
				UserID:   expenseModel.UserID,
				Amount:   expenseModel.Amount,
				Category: category,
			}

			expenses = append(expenses, expense)
		}
	}

	return expenses, nil
}

func (e *ExpenseRepository) GetExpense(expenseID string) (entities.Expense, error) {
	var expenseModel Expenses
	var categoryModel Categories

	result := e.gorm.Model(&Expenses{}).Where("id = ?", expenseID).First(&expenseModel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return entities.Expense{}, errors.New("expense not found")
		}
		return entities.Expense{}, errors.New(result.Error.Error())
	}

	result = e.gorm.Model(&Categories{}).Where("id = ?", expenseModel.CategoryID).First(&categoryModel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return entities.Expense{}, errors.New("error searching expense (" + expenseModel.ID + ") category")
		}
		return entities.Expense{}, errors.New(result.Error.Error())
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

	expense := entities.Expense{
		SharedEntity: entities.SharedEntity{
			ID:            expenseModel.ID,
			Active:        expenseModel.Active,
			CreatedAt:     expenseModel.CreatedAt,
			UpdatedAt:     expenseModel.UpdatedAt,
			DeactivatedAt: expenseModel.DeactivatedAt,
		},
		UserID:   expenseModel.UserID,
		Amount:   expenseModel.Amount,
		Category: category,
	}

	return expense, nil
}

func (e *ExpenseRepository) UpdateExpense(expense entities.Expense) error {
	result := e.gorm.Model(&Expenses{}).Where("id", expense.ID).Updates(Expenses{
		Amount:     expense.Amount,
		Notes:      expense.Notes,
		CategoryID: expense.Category.ID,
		UpdatedAt:  expense.UpdatedAt,
	})

	if result.Error != nil {
		return errors.New(result.Error.Error())
	}

	return nil
}
