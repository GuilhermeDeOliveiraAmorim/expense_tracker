package handlers

import (
	"net/http"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/infra/factory"
	usecases "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/use_cases"
	"github.com/gin-gonic/gin"
)

type ExpenseHandler struct {
	expenseFactory *factory.ExpenseFactory
}

func NewExpenseHandler(factory *factory.ExpenseFactory) *ExpenseHandler {
	return &ExpenseHandler{
		expenseFactory: factory,
	}
}

func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
	var input usecases.CreateExpenseInputDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, erros := h.expenseFactory.CreateExpense.Execute(input)
	if len(erros) > 0 {
		for _, err := range erros {
			if err.Status == 500 {
				c.JSON(err.Status, gin.H{"error": err})
				return
			} else {
				c.JSON(err.Status, gin.H{"error": err})
				return
			}
		}

	}

	c.JSON(http.StatusCreated, output)
}

func (h *ExpenseHandler) GetExpense(c *gin.Context) {
	var input usecases.GetExpenseInputDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.expenseFactory.GetExpense.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *ExpenseHandler) GetExpenses(c *gin.Context) {
	var input usecases.GetExpensesInputDto
	output, err := h.expenseFactory.GetExpenses.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *ExpenseHandler) UpdateExpense(c *gin.Context) {
	var input usecases.UpdateExpenseInputDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.expenseFactory.UpdateExpense.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *ExpenseHandler) DeleteExpense(c *gin.Context) {
	var input usecases.DeleteExpenseInputDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.expenseFactory.DeleteExpense.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, output)
}
