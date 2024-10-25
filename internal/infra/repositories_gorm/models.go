package repositoriesgorm

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Categories struct {
	ID            string    `gorm:"primaryKey;not null"`
	Active        bool      `gorm:"not null"`
	CreatedAt     time.Time `gorm:"not null"`
	UpdatedAt     time.Time `gorm:"not null"`
	DeactivatedAt time.Time `gorm:"not null"`
	UserID        string    `gorm:"not null"`
	Name          string    `gorm:"not null"`
	Color         string    `gorm:"not null"`
	User          Users     `gorm:"foreignKey:UserID"`
}

type Expenses struct {
	ID            string     `gorm:"primaryKey;not null"`
	Active        bool       `gorm:"not null"`
	CreatedAt     time.Time  `gorm:"not null"`
	UpdatedAt     time.Time  `gorm:"not null"`
	DeactivatedAt time.Time  `gorm:"not null"`
	UserID        string     `gorm:"not null"`
	Amount        float64    `gorm:"not null"`
	ExpanseDate   time.Time  `gorm:"not null"`
	CategoryID    string     `gorm:"not null"`
	Notes         string     `gorm:"null"`
	Category      Categories `gorm:"foreignKey:CategoryID"`
	Tags          []Tags     `gorm:"many2many:expense_tags"`
	User          Users      `gorm:"foreignKey:UserID"`
}

type Tags struct {
	ID            string    `gorm:"primaryKey;not null"`
	Active        bool      `gorm:"not null"`
	CreatedAt     time.Time `gorm:"not null"`
	UpdatedAt     time.Time `gorm:"not null"`
	DeactivatedAt time.Time `gorm:"not null"`
	UserID        string    `gorm:"not null"`
	Name          string    `gorm:"not null"`
	Color         string    `gorm:"not null"`
	User          Users     `gorm:"foreignKey:UserID"`
}

type Users struct {
	ID            string    `gorm:"primaryKey;not null"`
	Active        bool      `gorm:"not null"`
	CreatedAt     time.Time `gorm:"not null"`
	UpdatedAt     time.Time `gorm:"not null"`
	DeactivatedAt time.Time `gorm:"not null"`
	Name          string    `gorm:"not null"`
	Email         string    `gorm:"not null"`
	Password      string    `gorm:"not null"`
}

func Migration(db *gorm.DB, sqlDB *sql.DB) {
	if err := db.AutoMigrate(
		Categories{},
		Tags{},
		Expenses{},
		Users{},
	); err != nil {
		fmt.Println("Error during migration:", err)
		return
	}
	fmt.Println("Successful migration")
}
