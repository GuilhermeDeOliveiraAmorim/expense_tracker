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

func (c *ExpenseRepository) CreateExpense(expense entities.Expense) []error {
	if err := c.gorm.Create(&Expenses{
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
