package main

import (
	"fmt"
	"time"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/config"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/infra/factory"
	repositoriesgorm "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/infra/repositories_gorm"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/interface/handlers"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
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
		repositoriesgorm.Tags{},
		repositoriesgorm.Expenses{},
		repositoriesgorm.Users{},
	); err != nil {
		fmt.Println("Error during migration:", err)
		return
	}
	fmt.Println("Successful migration")

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
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

	public := r.Group("/")
	{
		public.POST("/users", userHandler.CreateUser)
		public.POST("/login", userHandler.Login)
	}

	protected := r.Group("/").Use(util.AuthMiddleware())
	{
		protected.POST("/categories", categoryHandler.CreateCategory)
		protected.GET("/categories", categoryHandler.GetCategory)
		protected.DELETE("/categories", categoryHandler.DeleteCategory)
		protected.GET("/categories/all", categoryHandler.GetCategories)
		protected.PATCH("/categories", categoryHandler.UpdateCategory)

		protected.POST("/tags", tagHandler.CreateTag)
		protected.GET("/tags", tagHandler.GetTag)
		protected.GET("/tags/all", tagHandler.GetTags)
		protected.PATCH("/tags", tagHandler.UpdateTag)
		protected.DELETE("/tags", tagHandler.DeleteTag)

		protected.POST("/expenses", expenseHandler.CreateExpense)
		protected.GET("/expenses", expenseHandler.GetExpense)
		protected.GET("/expenses/all", expenseHandler.GetExpenses)
		protected.PATCH("/expenses", expenseHandler.UpdateExpense)
		protected.DELETE("/expenses", expenseHandler.DeleteExpense)

		protected.GET("/users", userHandler.GetUser)
		protected.GET("/users/all", userHandler.GetUsers)
		protected.PATCH("/users", userHandler.UpdateUser)
		protected.DELETE("/users", userHandler.DeleteUser)

		protected.GET("/expenses/total", presentersHandler.GetTotalExpensesForPeriod)
		protected.GET("/expenses/categories", presentersHandler.GetExpensesByCategoryPeriod)
		protected.GET("/expenses/monthly/categories", presentersHandler.GetMonthlyExpensesByCategoryPeriod)
	}

	r.Run(":8080")
}
