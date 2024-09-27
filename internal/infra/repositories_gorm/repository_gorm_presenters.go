package repositoriesgorm

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type PresentersRepository struct {
	gorm *gorm.DB
}

func NewPresentersRepository(gorm *gorm.DB) *PresentersRepository {
	return &PresentersRepository{
		gorm: gorm,
	}
}

func (p *PresentersRepository) GetTotalExpensesForPeriod(userID string, startDate time.Time, endDate time.Time) (float64, error) {
	tx := p.gorm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var total float64

	if err := tx.Table("expenses").
		Select("SUM(amount)").
		Where("user_id = ? AND expanse_date BETWEEN ? AND ? AND active = ?", userID, startDate, endDate, true).
		Scan(&total).Error; err != nil {
		tx.Rollback()
		return 0, errors.New("failed to calculate total expenses: " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return 0, errors.New("failed to commit transaction")
	}

	return total, nil
}
