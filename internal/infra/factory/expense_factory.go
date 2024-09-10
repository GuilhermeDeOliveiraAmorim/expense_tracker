package factory

import (
	repositoriesgorm "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/infra/repositories_gorm"
	usecases "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/use_cases"
	"gorm.io/gorm"
)

type ExpenseFactory struct {
	CreateExpense *usecases.CreateExpenseUseCase
	DeleteExpense *usecases.DeleteExpenseUseCase
	GetExpenses   *usecases.GetExpensesUseCase
	GetExpense    *usecases.GetExpenseUseCase
	UpdateExpense *usecases.UpdateExpenseUseCase
}

func NewExpenseFactory(db *gorm.DB) *ExpenseFactory {
	expenseRepository := repositoriesgorm.NewExpenseRepository(db)
	userRepository := repositoriesgorm.NewUserRepository(db)

	createExpense := usecases.NewCreateExpenseUseCase(expenseRepository, userRepository)
	deleteExpense := usecases.NewDeleteExpenseUseCase(expenseRepository, userRepository)
	getExpenses := usecases.NewGetExpensesUseCase(expenseRepository, userRepository)
	getExpense := usecases.NewGetExpenseUseCase(expenseRepository, userRepository)
	updateExpense := usecases.NewUpdateExpenseUseCase(expenseRepository, userRepository)

	return &ExpenseFactory{
		CreateExpense: createExpense,
		DeleteExpense: deleteExpense,
		GetExpenses:   getExpenses,
		GetExpense:    getExpense,
		UpdateExpense: updateExpense,
	}
}
