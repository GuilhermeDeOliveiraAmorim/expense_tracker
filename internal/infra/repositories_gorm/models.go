package repositoriesgorm

import (
	"time"
)

type Categories struct {
	ID            string    `gorm:"primaryKey;not null"`
	Active        bool      `gorm:"not null"`
	CreatedAt     time.Time `gorm:"not null"`
	UpdatedAt     time.Time `gorm:"not null"`
	DeactivatedAt time.Time `gorm:"not null"`
	Name          string    `gorm:"not null"`
	Color         string    `gorm:"not null"`
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
	User          Users      `gorm:"foreignKey:UserID"`
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
