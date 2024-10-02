package handlers

import (
	"net/http"
	"strconv"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/infra/factory"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/presenters"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
	"github.com/gin-gonic/gin"
)

type PresentersHandler struct {
	presenterFactory *factory.PresentersFactory
}

func NewPresentersHandler(factory *factory.PresentersFactory) *PresentersHandler {
	return &PresentersHandler{
		presenterFactory: factory,
	}
}

func (h *PresentersHandler) GetTotalExpensesForPeriod(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
		return
	}

	startDate := c.Query("start_date")
	if startDate == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": util.ProblemDetails{
			Type:     "Bad Request",
			Title:    "Missing start date",
			Status:   http.StatusBadRequest,
			Detail:   "Start date is required",
			Instance: util.RFC400,
		}})
		return
	}

	endDate := c.Query("end_date")
	if endDate == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": util.ProblemDetails{
			Type:     "Bad Request",
			Title:    "Missing end date",
			Status:   http.StatusBadRequest,
			Detail:   "End date is required",
			Instance: util.RFC400,
		}})
		return
	}

	input := presenters.GetTotalExpensesForPeriodInputDto{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	output, errs := h.presenterFactory.GetTotalExpensesForPeriod.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *PresentersHandler) GetExpensesByCategoryPeriod(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
		return
	}

	startDate := c.Query("start_date")
	if startDate == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": util.ProblemDetails{
			Type:     "Bad Request",
			Title:    "Missing start date",
			Status:   http.StatusBadRequest,
			Detail:   "Start date is required",
			Instance: util.RFC400,
		}})
		return
	}

	endDate := c.Query("end_date")
	if endDate == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": util.ProblemDetails{
			Type:     "Bad Request",
			Title:    "Missing end date",
			Status:   http.StatusBadRequest,
			Detail:   "End date is required",
			Instance: util.RFC400,
		}})
		return
	}

	input := presenters.GetExpensesByCategoryPeriodInputDto{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	output, errs := h.presenterFactory.GetExpensesByCategoryPeriod.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *PresentersHandler) GetMonthlyExpensesByCategoryYear(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
		return
	}

	year := c.Query("year")
	if year == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": util.ProblemDetails{
			Type:     "Bad Request",
			Title:    "Missing year date",
			Status:   http.StatusBadRequest,
			Detail:   "Year date is required",
			Instance: util.RFC400,
		}})
		return
	}

	yearNumber, errYearNumber := strconv.Atoi(year)
	if errYearNumber != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": util.ProblemDetails{
			Type:     "Bad Request",
			Title:    "Invalid year date",
			Status:   http.StatusBadRequest,
			Detail:   "Year date is invalid",
			Instance: util.RFC400,
		}})
		return
	}

	input := presenters.GetMonthlyExpensesByCategoryYearInputDto{
		UserID: userID,
		Year:   yearNumber,
	}

	output, errs := h.presenterFactory.GetMonthlyExpensesByCategoryYear.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *PresentersHandler) GetMonthlyExpensesByTagYear(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
		return
	}

	year := c.Query("year")
	if year == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": util.ProblemDetails{
			Type:     "Bad Request",
			Title:    "Missing year date",
			Status:   http.StatusBadRequest,
			Detail:   "Year date is required",
			Instance: util.RFC400,
		}})
		return
	}

	yearNumber, errYearNumber := strconv.Atoi(year)
	if errYearNumber != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": util.ProblemDetails{
			Type:     "Bad Request",
			Title:    "Invalid year date",
			Status:   http.StatusBadRequest,
			Detail:   "Year date is invalid",
			Instance: util.RFC400,
		}})
		return
	}

	input := presenters.GetMonthlyExpensesByTagYearInputDto{
		UserID: userID,
		Year:   yearNumber,
	}

	output, errs := h.presenterFactory.GetMonthlyExpensesByTagYear.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}
