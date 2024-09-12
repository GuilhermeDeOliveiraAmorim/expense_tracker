package main

import (
	"fmt"
	"time"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/config"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/infra/factory"
	repositoriesgorm "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/infra/repositories_gorm"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/interface/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=" + config.DB_LOCAL.DB_HOST + " user=" + config.DB_LOCAL.DB_USER + " password=" + config.DB_LOCAL.DB_PASSWORD + " dbname=" + config.DB_LOCAL.DB_NAME + " port=" + config.DB_LOCAL.DB_PORT + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}
	fmt.Println("Successful connection")

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

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Permitir seu frontend acessar
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	categoryFactory := factory.NewCategoryFactory(db)
	categoryHandler := handlers.NewCategoryHandler(categoryFactory)

	tagFactory := factory.NewTagFactory(db)
	tagHandler := handlers.NewTagHandler(tagFactory)

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
	r.GET("/categories/category-treemap-amount-period", presentersHandler.ShowCategoryTreemapAmountPeriod)

	r.POST("/tags", tagHandler.CreateTag)
	r.GET("/tags", tagHandler.GetTag)
	r.DELETE("/tags", tagHandler.DeleteTag)
	r.GET("/tags/all", tagHandler.GetTags)

	r.POST("/expenses", expenseHandler.CreateExpense)
	r.GET("/expenses/all", expenseHandler.GetExpenses)
	r.GET("/expenses", expenseHandler.GetExpense)
	r.DELETE("/expenses", expenseHandler.DeleteExpense)
	r.PATCH("/expenses", expenseHandler.UpdateExpense)
	r.GET("/expenses/total-expenses-category-period", presentersHandler.ShowTotalExpensesCategoryPeriod)
	r.GET("/expenses/expense-simple-table-period", presentersHandler.ShowExpenseSimpleTablePeriod)

	r.POST("/users", userHandler.CreateUser)
	r.GET("/users/all", userHandler.GetUsers)
	r.GET("/users", userHandler.GetUser)
	r.DELETE("/users", userHandler.DeleteUser)
	r.PATCH("/users", userHandler.UpdateUser)
	r.POST("/login", userHandler.Login)

	r.Run(":8080")
}
