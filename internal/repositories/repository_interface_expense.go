package repositories

import "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"

type ExpenseRepositoryInterface interface {
	CreateExpense(expense entities.Expense) error
	DeleteExpense(expense entities.Expense) error
	GetExpenses() ([]entities.Expense, error)
	GetExpense(expenseID string) (entities.Expense, error)
	UpdateExpense(expense entities.Expense) error
}
