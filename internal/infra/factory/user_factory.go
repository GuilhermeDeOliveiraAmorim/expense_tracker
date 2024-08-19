package factory

import (
	repositoriesgorm "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/infra/repositories_gorm"
	usecases "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/use_cases"
	"gorm.io/gorm"
)

type UserFactory struct {
	CreateUser *usecases.CreateUserUseCase
	DeleteUser *usecases.DeleteUserUseCase
	GetUsers   *usecases.GetUsersUseCase
	GetUser    *usecases.GetUserUseCase
	UpdateUser *usecases.UpdateUserUseCase
}

func NewUserFactory(db *gorm.DB) *UserFactory {
	expenseRepository := repositoriesgorm.NewUserRepository(db)

	createUser := usecases.NewCreateUserUseCase(expenseRepository)
	deleteUser := usecases.NewDeleteUserUseCase(expenseRepository)
	getUsers := usecases.NewGetUsersUseCase(expenseRepository)
	getUser := usecases.NewGetUserUseCase(expenseRepository)
	updateUser := usecases.NewUpdateUserUseCase(expenseRepository)

	return &UserFactory{
		CreateUser: createUser,
		DeleteUser: deleteUser,
		GetUsers:   getUsers,
		GetUser:    getUser,
		UpdateUser: updateUser,
	}
}
