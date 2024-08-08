package repositories

import "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"

type UserRepositoryInterface interface {
	CreateUser(User entities.User) []error
	DeleteUser(User entities.User) []error
	GetUsers() ([]entities.User, []error)
	GetUser(userID string) (entities.User, []error)
	UpdateUser(User entities.User) []error
}
