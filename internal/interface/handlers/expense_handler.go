package handlers

import (
	"net/http"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/infra/factory"
	usecases "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/use_cases"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
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
	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
		return
	}

	var request CreateExpenseRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": util.ProblemDetails{
			Type:     "Bad Request",
			Title:    "Did not bind JSON",
			Status:   http.StatusBadRequest,
			Detail:   err.Error(),
			Instance: util.RFC400,
		}})
		return
	}

	input := usecases.CreateExpenseInputDto{
		UserID:      userID,
		CategoryID:  request.CategoryID,
		Amount:      request.Amount,
		Tags:        request.Tags,
		ExpenseDate: request.ExpenseDate,
		Notes:       request.Notes,
	}

	output, errs := h.expenseFactory.CreateExpense.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusCreated, output)
}

func (h *ExpenseHandler) GetExpense(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
		return
	}

	expenseID := c.Query("expense_id")
	if expenseID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": util.ProblemDetails{
			Type:     "Bad Request",
			Title:    "Missing Expense ID",
			Status:   http.StatusBadRequest,
			Detail:   "Expense id is required",
			Instance: util.RFC400,
		}})
		return
	}

	input := usecases.GetExpenseInputDto{
		UserID:    userID,
		ExpenseID: expenseID,
	}

	output, errs := h.expenseFactory.GetExpense.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *ExpenseHandler) GetExpenses(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
		return
	}

	input := usecases.GetExpensesInputDto{
		UserID: userID,
	}

	output, errs := h.expenseFactory.GetExpenses.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *ExpenseHandler) UpdateExpense(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
		return
	}

	var request UpdateExpenseRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": util.ProblemDetails{
			Type:     "Bad Request",
			Title:    "Did not bind JSON",
			Status:   http.StatusBadRequest,
			Detail:   err.Error(),
			Instance: util.RFC400,
		}})
		return
	}

	input := usecases.UpdateExpenseInputDto{
		UserID:      userID,
		ExpenseID:   request.ExpenseID,
		Amount:      request.Amount,
		ExpenseDate: request.ExpenseDate,
		CategoryID:  request.CategoryID,
		Notes:       request.Notes,
	}

	output, erros := h.expenseFactory.UpdateExpense.Execute(input)
	if len(erros) > 0 {
		handleErrors(c, erros)
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *ExpenseHandler) DeleteExpense(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
		return
	}

	expenseID := c.Query("expense_id")
	if expenseID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": util.ProblemDetails{
			Type:     "Bad Request",
			Title:    "Missing Expense ID",
			Status:   http.StatusBadRequest,
			Detail:   "Expense id is required",
			Instance: util.RFC400,
		}})
		return
	}

	input := usecases.DeleteExpenseInputDto{
		UserID:    userID,
		ExpenseID: expenseID,
	}

	output, errs := h.expenseFactory.DeleteExpense.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}
