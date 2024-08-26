package repositories

import "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"

type UserRepositoryInterface interface {
	CreateUser(user entities.User) error
	DeleteUser(user entities.User) error
	GetUsers() ([]entities.User, error)
	GetUser(userID string) (entities.User, error)
	UpdateUser(user entities.User) error
}
