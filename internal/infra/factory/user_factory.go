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
	Login      *usecases.LoginUseCase
}

func NewUserFactory(db *gorm.DB) *UserFactory {
	userRepository := repositoriesgorm.NewUserRepository(db)

	createUser := usecases.NewCreateUserUseCase(userRepository)
	deleteUser := usecases.NewDeleteUserUseCase(userRepository)
	getUsers := usecases.NewGetUsersUseCase(userRepository)
	getUser := usecases.NewGetUserUseCase(userRepository)
	updateUser := usecases.NewUpdateUserUseCase(userRepository)
	login := usecases.NewLoginUseCase(userRepository)

	return &UserFactory{
		CreateUser: createUser,
		DeleteUser: deleteUser,
		GetUsers:   getUsers,
		GetUser:    getUser,
		UpdateUser: updateUser,
		Login:      login,
	}
}
