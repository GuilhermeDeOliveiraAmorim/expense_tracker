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
	tx := u.gorm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := tx.Create(&Users{
		ID:            user.ID,
		Active:        user.Active,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		DeactivatedAt: user.DeactivatedAt,
		Name:          user.Name,
		Email:         user.Login.Email,
		Password:      user.Login.Password,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return errors.New("failed to commit transaction: " + err.Error())
	}

	return nil
}

func (u *UserRepository) DeleteUser(user entities.User) error {
	tx := u.gorm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	result := tx.Model(&Users{}).Where("id = ? AND active = ?", user.ID, true).
		Select("Active", "DeactivatedAt", "UpdatedAt").Updates(Users{
		Active:        user.Active,
		DeactivatedAt: user.DeactivatedAt,
		UpdatedAt:     user.UpdatedAt,
	})

	if result.Error != nil {
		tx.Rollback()
		return errors.New(result.Error.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return errors.New("failed to commit transaction: " + err.Error())
	}

	return nil
}

func (u *UserRepository) GetUsers() ([]entities.User, error) {
	var usersModel []Users
	if err := u.gorm.Find(&usersModel).Error; err != nil {
		u.gorm.Rollback()
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
	} else {
		users = []entities.User{}
	}

	return users, nil
}

func (u *UserRepository) GetUser(userID string) (entities.User, error) {
var userModel Users

	result := u.gorm.Model(&Users{}).Where("id = ?", userID).First(&userModel)
	if result.Error != nil {
		u.gorm.Rollback()
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

func (u *UserRepository) UpdateUser(user entities.User) error {
	tx := u.gorm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	result := tx.Model(&Users{}).Where("id", user.ID).Updates(Users{
		Name:      user.Name,
		UpdatedAt: user.UpdatedAt,
	})

	if result.Error != nil {
		tx.Rollback()
		return errors.New(result.Error.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return errors.New("failed to commit transaction: " + err.Error())
	}

	return nil
}

func (u *UserRepository) GetUserByEmail(userEmail string) (entities.User, error) {
	tx := u.gorm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var userModel Users

	result := tx.Model(&Users{}).Where("email = ?", userEmail).First(&userModel)
	if result.Error != nil {
		tx.Rollback()
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
		Login: entities.Login{
			Email:    userModel.Email,
			Password: userModel.Password,
		},
	}

	return user, nil
}

func (u *UserRepository) ThisUserExists(userName string) (bool, error) {
	tx := u.gorm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var userModel Users

	result := tx.Model(&Users{}).Where("name = ?", userName).First(&userModel)
	if result.Error != nil {
		tx.Rollback()
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, errors.New("user not found")
		}
		return false, errors.New(result.Error.Error())
	}

	return true, nil
}

func (u *UserRepository) ThisUserEmailExists(userEmail string) (bool, error) {
	tx := u.gorm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var userModel Users

	result := tx.Model(&Users{}).Where("email = ?", userEmail).First(&userModel)
	if result.Error != nil {
		tx.Rollback()
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, errors.New("not found")
		}
		return false, errors.New(result.Error.Error())
	}

	return true, nil
}

func (u *UserRepository) ThisUserNameExists(userName string) (bool, error) {
	tx := u.gorm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var userModel Users

	result := tx.Model(&Users{}).Where("name = ?", userName).First(&userModel)
	if result.Error != nil {
		tx.Rollback()
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, errors.New("not found")
		}
		return false, errors.New(result.Error.Error())
	}

	return true, nil
}
