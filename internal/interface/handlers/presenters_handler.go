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

// @Summary Get total expenses for a period
// @Description Retrieves the total expenses of a user for a specified date range
// @Tags Presenters
// @Produce json
// @Param start_date query string true "Start date (DDMMYYYY)"
// @Param end_date query string true "End date (DDMMYYYY)"
// @Success 200 {object} presenters.GetTotalExpensesForPeriodOutputDto
// @Failure 400 {object} util.ProblemDetails "Bad Request - Missing start date or end date"
// @Failure 401 {object} util.ProblemDetails "Unauthorized"
// @Failure 500 {object} util.ProblemDetails "Internal Server Error"
// @Security BearerAuth
// @Router /expenses/total [get]
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

// @Summary Get expenses by category for a period
// @Description Retrieves the expenses of a user categorized by category for a specified date range
// @Tags Presenters
// @Produce json
// @Param start_date query string true "Start date (DDMMYYYY)"
// @Param end_date query string true "End date (DDMMYYYY)"
// @Success 200 {object} presenters.GetExpensesByCategoryPeriodOutputDto
// @Failure 400 {object} util.ProblemDetails "Bad Request - Missing start date or end date"
// @Failure 401 {object} util.ProblemDetails "Unauthorized"
// @Failure 500 {object} util.ProblemDetails "Internal Server Error"
// @Security BearerAuth
// @Router /expenses/categories [get]
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

// @Summary Get monthly expenses by category for a specific year
// @Description Retrieves the monthly expenses of a user categorized by category for a specific year
// @Tags Presenters
// @Produce json
// @Param year query string true "Year (YYYY)"
// @Success 200 {object} presenters.GetMonthlyExpensesByCategoryYearOutputDto
// @Failure 400 {object} util.ProblemDetails "Bad Request - Missing or invalid year"
// @Failure 401 {object} util.ProblemDetails "Unauthorized"
// @Failure 500 {object} util.ProblemDetails "Internal Server Error"
// @Security BearerAuth
// @Router /expenses/categories/monthly [get]
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

	input := presenters.GetMonthlyExpensesByCategoryYearInputDto{
		UserID: userID,
		Year:   year,
	}

	output, errs := h.presenterFactory.GetMonthlyExpensesByCategoryYear.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

// @Summary Get monthly expenses by tag for a specific year
// @Description Retrieves the monthly expenses of a user categorized by tags for a specific year
// @Tags Presenters
// @Produce json
// @Param year query string true "Year (YYYY)"
// @Success 200 {object} presenters.GetMonthlyExpensesByTagYearOutputDto
// @Failure 400 {object} util.ProblemDetails "Bad Request - Missing or invalid year"
// @Failure 401 {object} util.ProblemDetails "Unauthorized"
// @Failure 500 {object} util.ProblemDetails "Internal Server Error"
// @Security BearerAuth
// @Router /expenses/tags/monthly [get]
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

	input := presenters.GetMonthlyExpensesByTagYearInputDto{
		UserID: userID,
		Year:   year,
	}

	output, errs := h.presenterFactory.GetMonthlyExpensesByTagYear.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

// @Summary Get total expenses for the current month
// @Description Retrieves the total expenses of a user for the current month
// @Tags Presenters
// @Produce json
// @Success 200 {object} presenters.GetTotalExpensesForCurrentMonthOutputDto
// @Failure 400 {object} util.ProblemDetails "Bad Request - Invalid parameters"
// @Failure 401 {object} util.ProblemDetails "Unauthorized"
// @Failure 500 {object} util.ProblemDetails "Internal Server Error"
// @Security BearerAuth
// @Router /expenses/monthly/total [get]
func (h *PresentersHandler) GetTotalExpensesForCurrentMonth(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
		return
	}

	input := presenters.GetTotalExpensesForCurrentMonthInputDto{
		UserID: userID,
	}

	output, errs := h.presenterFactory.GetTotalExpensesForCurrentMonth.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

// @Summary Get expenses by month and year
// @Description Retrieves expenses of a user for a specific month and year
// @Tags Presenters
// @Produce json
// @Param year query string true "Year of the expenses" Format(year)
// @Param month query string true "Month of the expenses" Format(month)
// @Success 200 {object} presenters.GetExpensesByMonthYearOutputDto
// @Failure 400 {object} util.ProblemDetails "Bad Request - Missing or invalid parameters"
// @Failure 401 {object} util.ProblemDetails "Unauthorized"
// @Failure 500 {object} util.ProblemDetails "Internal Server Error"
// @Security BearerAuth
// @Router /expenses/monthly/year [get]
func (h *PresentersHandler) GetExpensesByMonthYear(c *gin.Context) {
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

	month := c.Query("month")
	if month == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": util.ProblemDetails{
			Type:     "Bad Request",
			Title:    "Missing month date",
			Status:   http.StatusBadRequest,
			Detail:   "Year date is required",
			Instance: util.RFC400,
		}})
		return
	}

	input := presenters.GetExpensesByMonthYearInputDto{
		UserID: userID,
		Year:   year,
		Month:  month,
	}

	output, errs := h.presenterFactory.GetExpensesByMonthYear.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

// @Summary Get total expenses for the current week
// @Description Retrieves the total expenses of a user for the current week
// @Tags Presenters
// @Produce json
// @Success 200 {object} presenters.GetTotalExpensesForCurrentWeekOutputDto
// @Failure 400 {object} util.ProblemDetails "Bad Request - Invalid parameters"
// @Failure 401 {object} util.ProblemDetails "Unauthorized"
// @Failure 500 {object} util.ProblemDetails "Internal Server Error"
// @Security BearerAuth
// @Router /expenses/weekly/total [get]
func (h *PresentersHandler) GetTotalExpensesForCurrentWeek(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
		return
	}

	input := presenters.GetTotalExpensesForCurrentWeekInputDto{
		UserID: userID,
	}

	output, errs := h.presenterFactory.GetTotalExpensesForCurrentWeek.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

// @Summary Get total expenses for each month in the specified year
// @Description Retrieves the total expenses for each month of the specified year for a user
// @Tags Presenters
// @Produce json
// @Param year query string true "Year for which to retrieve monthly expenses"
// @Success 200 {object} presenters.GetTotalExpensesMonthCurrentYearOutputDto
// @Failure 400 {object} util.ProblemDetails "Bad Request - Missing or invalid year"
// @Failure 401 {object} util.ProblemDetails "Unauthorized"
// @Failure 500 {object} util.ProblemDetails "Internal Server Error"
// @Security BearerAuth
// @Router /expenses/total/monthly/year [get]
func (h *PresentersHandler) GetTotalExpensesMonthCurrentYear(c *gin.Context) {
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

	input := presenters.GetTotalExpensesMonthCurrentYearInputDto{
		UserID: userID,
		Year:   year,
	}

	output, errs := h.presenterFactory.GetTotalExpensesMonthCurrentYear.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

// @Summary Get totals of expenses by category tags for a specific month and year
// @Description Retrieves the total expenses by category tags for a given month and year for a user
// @Tags Presenters
// @Produce json
// @Param year query string true "Year for which to retrieve expenses by category tags"
// @Param month query string true "Month for which to retrieve expenses by category tags"
// @Success 200 {object} presenters.GetCategoryTagsTotalsByMonthYearOutputDto
// @Failure 400 {object} util.ProblemDetails "Bad Request - Missing or invalid year/month"
// @Failure 401 {object} util.ProblemDetails "Unauthorized"
// @Failure 500 {object} util.ProblemDetails "Internal Server Error"
// @Security BearerAuth
// @Router /expenses/tags/monthly/total [get]
func (h *PresentersHandler) GetCategoryTagsTotalsByMonthYear(c *gin.Context) {
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

	month := c.Query("month")
	if month == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": util.ProblemDetails{
			Type:     "Bad Request",
			Title:    "Missing month date",
			Status:   http.StatusBadRequest,
			Detail:   "Year date is required",
			Instance: util.RFC400,
		}})
		return
	}

	input := presenters.GetCategoryTagsTotalsByMonthYearInputDto{
		UserID: userID,
		Year:   year,
		Month:  month,
	}

	output, errs := h.presenterFactory.GetCategoryTagsTotalsByMonthYear.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

// @Summary Get available months and years
// @Description Retrieves the list of months and years for which expense data is available for a user
// @Tags Utility
// @Produce json
// @Success 200 {object} presenters.GetAvailableMonthsYearsOutputDto
// @Failure 401 {object} util.ProblemDetails "Unauthorized"
// @Failure 500 {object} util.ProblemDetails "Internal Server Error"
// @Security BearerAuth
// @Router /expenses/available-months-years [get]
func (h *PresentersHandler) GetAvailableMonthsYears(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
		return
	}

	input := presenters.GetAvailableMonthsYearsInputDto{
		UserID: userID,
	}

	output, errs := h.presenterFactory.GetAvailableMonthsYears.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

// @Summary Get expenses in a time interval
// @Description A date range is passed to return a set of expenses in that range
// @Tags Presenters
// @Produce json
// @Success 200 {object} presenters.GetDayToDayExpensesPeriodOutputDto
// @Failure 401 {object} util.ProblemDetails "Unauthorized"
// @Failure 500 {object} util.ProblemDetails "Internal Server Error"
// @Security BearerAuth
// @Router /expenses/day/day/period [get]
func (h *PresentersHandler) GetDayToDayExpensesPeriod(c *gin.Context) {
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

	input := presenters.GetDayToDayExpensesPeriodInputDto{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	output, errs := h.presenterFactory.GetDayToDayExpensesPeriod.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}
