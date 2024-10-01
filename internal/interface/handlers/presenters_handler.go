package handlers

import (
	"net/http"

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

func (h *PresentersHandler) GetMonthlyExpensesByCategoryPeriod(c *gin.Context) {
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

	input := presenters.GetMonthlyExpensesByCategoryPeriodInputDto{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	output, errs := h.presenterFactory.GetMonthlyExpensesByCategoryPeriod.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}
