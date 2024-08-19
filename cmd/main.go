package main

import (
	"fmt"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/configs"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/infra/factory"
	repositoriesgorm "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/infra/repositories_gorm"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/interface/handlers"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	dsn := "host=" + configs.PostgresServer + " user=" + configs.PostgresUser + " password=" + configs.PostgresPassword + " dbname=" + configs.PostgresDb + " port=" + configs.PostgresPort + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err := db.AutoMigrate(
		repositoriesgorm.Categories{},
		repositoriesgorm.Expenses{},
		repositoriesgorm.Users{},
	); err != nil {
		fmt.Println("Erro durante a migração:", err)
		return
	}
	fmt.Println("Migração bem-sucedida!")

	r := gin.Default()

	categoryFactory := factory.NewCategoryFactory(db)
	categoryHandler := handlers.NewCategoryHandler(categoryFactory)

	r.POST("/categories", categoryHandler.CreateCategory)
	r.GET("/categories/:category_id/categories", categoryHandler.GetCategory)

	r.Run(":8080")
}
