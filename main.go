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
)

func main() {
	db, sqlDB, err := util.SetupDatabaseConnection(util.NEON)
	if err != nil {
		panic("Failed to connect database")
	}
	fmt.Println("Successful connection")

	repositoriesgorm.Migration(db, sqlDB)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.FRONT_END_URL_VAR.FRONT_END_URL_DEV, config.FRONT_END_URL_VAR.FRONT_END_URL_PROD},
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
		protected.GET("/expenses/categories/monthly", presentersHandler.GetMonthlyExpensesByCategoryYear)
		protected.GET("/expenses/tags/monthly", presentersHandler.GetMonthlyExpensesByTagYear)
		protected.GET("/expenses/monthly/total", presentersHandler.GetTotalExpensesForCurrentMonth)
		protected.GET("/expenses/monthly/year", presentersHandler.GetExpensesByMonthYear)
		protected.GET("/expenses/weekly/total", presentersHandler.GetTotalExpensesForCurrentWeek)
		protected.GET("/expenses/total/monthly/year", presentersHandler.GetTotalExpensesMonthCurrentYear)
		protected.GET("/expenses/tags/monthly/total", presentersHandler.GetCategoryTagsTotalsByMonthYear)

		protected.GET("/util/months/years", presentersHandler.GetAvailableMonthsYears)
	}

	r.Run(":8080")
}
