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

// @Summary      Create an expense
// @Description  Create a new expense entry
// @Tags         Expenses
// @Accept       json
// @Produce      json
// @Param        request body CreateExpenseRequest true "Expense data"
// @Success      201 {object} usecases.CreateExpenseOutputDto
// @Failure      400 {object} util.ProblemDetails "Bad Request"
// @Failure      500 {object} util.ProblemDetails "Internal Server Error"
// @Failure		 401 {object} util.ProblemDetails "Unauthorized"
// @Security	 BearerAuth
// @Router       /expenses [post]
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

// @Summary      Get a specific expense
// @Description  Retrieve an expense by its ID
// @Tags         Expenses
// @Accept       json
// @Produce      json
// @Param        expense_id query string true "Expense ID"
// @Success      200 {object} usecases.GetExpenseOutputDto
// @Failure      400 {object} util.ProblemDetails "Bad Request"
// @Failure      401 {object} util.ProblemDetails "Unauthorized"
// @Failure		 401 {object} util.ProblemDetails "Unauthorized"
// @Security	 BearerAuth
// @Failure      500 {object} util.ProblemDetails "Internal Server Error"
// @Router       /expenses/{expense_id} [get]
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

// @Summary      Get all expenses
// @Description  Retrieve all expenses for the authenticated user
// @Tags         Expenses
// @Accept       json
// @Produce      json
// @Success      200 {array} usecases.GetExpensesOutputDto
// @Failure		 401 {object} util.ProblemDetails "Unauthorized"
// @Security	 BearerAuth
// @Failure      500 {object} util.ProblemDetails "Internal Server Error"
// @Router       /expenses/all [get]
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

// @Summary      Update an expense
// @Description  Update an existing expense
// @Tags         Expenses
// @Accept       json
// @Produce      json
// @Param        request body UpdateExpenseRequest true "Updated expense data"
// @Success      200 {object} usecases.UpdateExpenseOutputDto
// @Failure      400 {object} util.ProblemDetails "Bad Request"
// @Failure		 401 {object} util.ProblemDetails "Unauthorized"
// @Security	 BearerAuth
// @Failure      404 {object} util.ProblemDetails "Expense Not Found"
// @Failure      500 {object} util.ProblemDetails "Internal Server Error"
// @Router       /expenses/{expense_id} [patch]
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
		Tags:        request.Tags,
	}

	output, erros := h.expenseFactory.UpdateExpense.Execute(input)
	if len(erros) > 0 {
		handleErrors(c, erros)
		return
	}

	c.JSON(http.StatusOK, output)
}

// @Summary      Delete an expense
// @Description  Delete an expense by its ID
// @Tags         Expenses
// @Accept       json
// @Produce      json
// @Param        expense_id query string true "Expense ID"
// @Success      200 {object} usecases.DeleteExpenseOutputDto
// @Failure      400 {object} util.ProblemDetails "Bad Request"
// @Failure		 401 {object} util.ProblemDetails "Unauthorized"
// @Security	 BearerAuth
// @Failure      404 {object} util.ProblemDetails "Expense Not Found"
// @Failure      500 {object} util.ProblemDetails "Internal Server Error"
// @Router       /expenses/{expense_id} [delete]
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
