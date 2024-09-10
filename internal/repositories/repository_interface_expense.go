package repositories

import "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"

type ExpenseRepositoryInterface interface {
	CreateExpense(expense entities.Expense) error
	DeleteExpense(expense entities.Expense) error
	GetExpenses(userID string) ([]entities.Expense, error)
	GetExpense(userID string, expenseID string) (entities.Expense, error)
	UpdateExpense(expense entities.Expense) error
}
