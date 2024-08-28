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
		panic("Failed to connect database")
	}

	if err := db.AutoMigrate(
		repositoriesgorm.Categories{},
		repositoriesgorm.Expenses{},
		repositoriesgorm.Users{},
	); err != nil {
		fmt.Println("Error during migration:", err)
		return
	}
	fmt.Println("Successful migration")

	r := gin.Default()

	categoryFactory := factory.NewCategoryFactory(db)
	categoryHandler := handlers.NewCategoryHandler(categoryFactory)

	expenseFactory := factory.NewExpenseFactory(db)
	expenseHandler := handlers.NewExpenseHandler(expenseFactory)

	userFactory := factory.NewUserFactory(db)
	userHandler := handlers.NewUserHandler(userFactory)

	presentersFactory := factory.NewPresentersFactory(db)
	presentersHandler := handlers.NewPresentersHandler(presentersFactory)

	r.POST("/categories", categoryHandler.CreateCategory)
	r.GET("/categories", categoryHandler.GetCategory)
	r.DELETE("/categories", categoryHandler.DeleteCategory)
	r.GET("/categories/all", categoryHandler.GetCategories)
	r.PATCH("/categories", categoryHandler.UpdateCategory)

	r.POST("/expenses", expenseHandler.CreateExpense)
	r.GET("/expenses/all", expenseHandler.GetExpenses)
	r.GET("/expenses", expenseHandler.GetExpense)
	r.DELETE("/expenses", expenseHandler.DeleteExpense)
	r.PATCH("/expenses", expenseHandler.UpdateExpense)

	r.POST("/users", userHandler.CreateUser)
	r.GET("/users/all", userHandler.GetUsers)
	r.GET("/users", userHandler.GetUser)
	r.DELETE("/users", userHandler.DeleteUser)
	r.PATCH("/users", userHandler.UpdateUser)

	r.GET("/users/total-expenses-category-period", presentersHandler.ShowTotalExpensesCategoryPeriod)

	r.Run(":8080")
}
