package repositories

import "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"

type GetExpense struct {
	Expense  entities.Expense  `json:"expense"`
	Category entities.Category `json:"category"`
	Tags     []entities.Tag    `json:"tags"`
}

type ExpenseRepositoryInterface interface {
	CreateExpense(expense entities.Expense) error
	DeleteExpense(expense entities.Expense) error
	GetExpenses(userID string) ([]GetExpense, error)
	GetExpense(userID string, expenseID string) (GetExpense, error)
	UpdateExpense(expense entities.Expense) error
}
