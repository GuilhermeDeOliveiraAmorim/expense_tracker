package repositoriesgorm

import (
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

func (e *ExpenseRepository) CreateExpense(expense entities.Expense) []error {
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
		return []error{err}
	}

	return nil
}

func (e *ExpenseRepository) DeleteExpense(category entities.Expense) []error {
	err := e.gorm.Model(&Expenses{}).Where("id = ?", category.ID).Updates(Expenses{
		Active:        category.Active,
		DeactivatedAt: category.DeactivatedAt,
	}).Error

	if err != nil {
		return []error{err}
	}

	return nil
}
