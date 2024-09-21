package repositoriesgorm

import (
	"errors"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"gorm.io/gorm"
)

type ExpenseRepository struct {
	gorm *gorm.DB
}

func NewExpenseRepository(gorm *gorm.DB) *ExpenseRepository {
	return &ExpenseRepository{
		gorm: gorm,
	}
}

func (e *ExpenseRepository) CreateExpense(expense entities.Expense) error {
	tx := e.gorm.Begin()

	if err := tx.Create(&Expenses{
		ID:            expense.ID,
		Active:        expense.Active,
		CreatedAt:     expense.CreatedAt,
		UpdatedAt:     expense.UpdatedAt,
		DeactivatedAt: expense.DeactivatedAt,
		UserID:        expense.UserID,
		Amount:        expense.Amount,
		ExpanseDate:   expense.ExpenseDate,
		CategoryID:    expense.CategoryID,
		Notes:         expense.Notes,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, tagID := range expense.TagIDs {
		if err := tx.Exec("INSERT INTO expense_tags (expenses_id, tags_id) VALUES (?, ?)", expense.ID, tagID).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (e *ExpenseRepository) DeleteExpense(expense entities.Expense) error {
	result := e.gorm.Model(&Expenses{}).Where("id = ? AND user_id = ? AND active = true", expense.ID, expense.UserID).
		Select("Active", "DeactivatedAt", "UpdatedAt").Updates(Expenses{
		Active:        expense.Active,
		DeactivatedAt: expense.DeactivatedAt,
		UpdatedAt:     expense.UpdatedAt,
	})

	if result.Error != nil {
		return errors.New(result.Error.Error())
	}

	return nil
}

func (e *ExpenseRepository) GetExpenses(userID string) ([]entities.Expense, error) {
	var expensesModel []Expenses

	if err := e.gorm.Preload("Tags").Preload("Category").Where("user_id = ? AND active = true", userID).Find(&expensesModel).Error; err != nil {
		return nil, err
	}

	var expenses []entities.Expense

	if len(expensesModel) > 0 {
		for _, expenseModel := range expensesModel {
			var tagIDs []string
			for _, tag := range expenseModel.Tags {
				tagIDs = append(tagIDs, tag.ID)
			}

			expense := entities.Expense{
				SharedEntity: entities.SharedEntity{
					ID:            expenseModel.ID,
					Active:        expenseModel.Active,
					CreatedAt:     expenseModel.CreatedAt,
					UpdatedAt:     expenseModel.UpdatedAt,
					DeactivatedAt: expenseModel.DeactivatedAt,
				},
				UserID:      expenseModel.UserID,
				Amount:      expenseModel.Amount,
				ExpenseDate: expenseModel.ExpanseDate,
				Notes:       expenseModel.Notes,
				CategoryID:  expenseModel.Category.ID,
				TagIDs:      tagIDs,
			}

			expenses = append(expenses, expense)
		}
	}

	return expenses, nil
}

func (e *ExpenseRepository) GetExpense(userID string, expenseID string) (entities.Expense, error) {
	var expenseModel Expenses

	result := e.gorm.Preload("Tags").Preload("Category").Where("id = ? AND user_id = ? AND active = true", expenseID, userID).First(&expenseModel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return entities.Expense{}, errors.New("expense not found")
		}
		return entities.Expense{}, errors.New(result.Error.Error())
	}

	var tagIDs []string

	for _, tag := range expenseModel.Tags {
		tagIDs = append(tagIDs, tag.ID)
	}

	expense := entities.Expense{
		SharedEntity: entities.SharedEntity{
			ID:            expenseModel.ID,
			Active:        expenseModel.Active,
			CreatedAt:     expenseModel.CreatedAt,
			UpdatedAt:     expenseModel.UpdatedAt,
			DeactivatedAt: expenseModel.DeactivatedAt,
		},
		UserID:      expenseModel.UserID,
		Amount:      expenseModel.Amount,
		ExpenseDate: expenseModel.ExpanseDate,
		Notes:       expenseModel.Notes,
		CategoryID:  expenseModel.Category.ID,
		TagIDs:      tagIDs,
	}

	return expense, nil
}

func (e *ExpenseRepository) UpdateExpense(expense entities.Expense) error {
	result := e.gorm.Model(&Expenses{}).Where("id AND active = true", expense.ID).Updates(Expenses{
		Amount:     expense.Amount,
		Notes:      expense.Notes,
		CategoryID: expense.CategoryID,
		UpdatedAt:  expense.UpdatedAt,
	})

	if result.Error != nil {
		return errors.New(result.Error.Error())
	}

	return nil
}
