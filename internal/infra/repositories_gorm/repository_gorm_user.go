package repositoriesgorm

import (
	"errors"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"gorm.io/gorm"
)

type UserRepository struct {
	gorm *gorm.DB
}

func NewUserRepository(gorm *gorm.DB) *UserRepository {
	return &UserRepository{
		gorm: gorm,
	}
}

func (u *UserRepository) CreateUser(user entities.User) error {
	if err := u.gorm.Create(&Users{
		ID:            user.ID,
		Active:        user.Active,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		DeactivatedAt: user.DeactivatedAt,
		Name:          user.Name,
		Email:         user.Login.Email,
		Password:      user.Login.Password,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (u *UserRepository) DeleteUser(user entities.User) error {
	err := u.gorm.Model(&Users{}).Where("id = ?", user.ID).Updates(Users{
		Active:        user.Active,
		DeactivatedAt: user.DeactivatedAt,
	}).Error

	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepository) GetUsers() ([]entities.User, error) {
	var usersModel []Users
	if err := u.gorm.Find(&usersModel).Error; err != nil {
		return nil, err
	}

	var users []entities.User

	if len(usersModel) > 0 {
		for _, userModel := range usersModel {
			user := entities.User{
				SharedEntity: entities.SharedEntity{
					ID:            userModel.ID,
					Active:        userModel.Active,
					CreatedAt:     userModel.CreatedAt,
					UpdatedAt:     userModel.UpdatedAt,
					DeactivatedAt: userModel.DeactivatedAt,
				},
				Name: userModel.Name,
			}

			users = append(users, user)
		}
	}

	return users, nil
}

func (u *UserRepository) GetUser(userID string) (entities.User, error) {
	var userModel Users

	result := u.gorm.Model(&Users{}).Where("id = ?", userID).First(&userModel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return entities.User{}, errors.New("user not found")
		}
		return entities.User{}, errors.New(result.Error.Error())
	}

	user := entities.User{
		SharedEntity: entities.SharedEntity{
			ID:            userModel.ID,
			Active:        userModel.Active,
			CreatedAt:     userModel.CreatedAt,
			UpdatedAt:     userModel.UpdatedAt,
			DeactivatedAt: userModel.DeactivatedAt,
		},
		Name: userModel.Name,
	}

	return user, nil
}

func (c *UserRepository) UpdateUser(user entities.User) error {
	result := c.gorm.Model(&Users{}).Where("id", user.ID).Updates(Users{
		Name:      user.Name,
		UpdatedAt: user.UpdatedAt,
	})

	if result.Error != nil {
		return errors.New(result.Error.Error())
	}

	return nil
}
