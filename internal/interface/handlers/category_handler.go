package handlers

import (
	"net/http"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/infra/factory"
	usecases "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/use_cases"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryFactory *factory.CategoryFactory
}

func NewCategoryHandler(factory *factory.CategoryFactory) *CategoryHandler {
	return &CategoryHandler{
		categoryFactory: factory,
	}
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new category with the provided details
// @Tags categories
// @Accept json
// @Produce json
// @Success 201 {object} usecases.CreateCategoryOutputDto
// @Failure 400 {object} util.ProblemDetails
// @Failure 500 {object} util.ProblemDetails
// @Param request body CreateCategoryRequest true "Request body to create a new category"
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
		return
	}

	var request CreateCategoryRequest
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

	input := usecases.CreateCategoryInputDto{
		UserID: userID,
		Name:   request.Name,
		Color:  request.Color,
	}

	output, errs := h.categoryFactory.CreateCategory.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusCreated, output)
}

// GetCategory godoc
// @Summary Get category details
// @Description Get details of a category by its ID
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {object} usecases.GetCategoryOutputDto
// @Failure 400 {object} util.ProblemDetails
// @Failure 404 {object} util.ProblemDetails
// @Failure 500 {object} util.ProblemDetails
// @Param category_id query string true "Category ID"
// @Router /categories/{category_id} [get]
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
		return
	}

	categoryID := c.Query("category_id")
	if categoryID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": util.ProblemDetails{
			Type:     "Bad Request",
			Title:    "Missing Category ID",
			Status:   http.StatusBadRequest,
			Detail:   "Category id is required",
			Instance: util.RFC400,
		}})
		return
	}

	input := usecases.GetCategoryInputDto{
		UserID:     userID,
		CategoryID: categoryID,
	}

	output, errs := h.categoryFactory.GetCategory.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

// GetCategories godoc
// @Summary Get all categories
// @Description Retrieve all categories for the authenticated user
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {array} usecases.GetCategoriesOutputDto
// @Failure 400 {object} util.ProblemDetails
// @Failure 500 {object} util.ProblemDetails
// @Router /categories [get]
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
		return
	}

	input := usecases.GetCategoriesInputDto{
		UserID: userID,
	}

	output, errs := h.categoryFactory.GetCategories.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

// UpdateCategory godoc
// @Summary Update a category
// @Description Update details of an existing category
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {object} usecases.UpdateCategoryOutputDto
// @Failure 400 {object} util.ProblemDetails
// @Failure 404 {object} util.ProblemDetails
// @Failure 500 {object} util.ProblemDetails
// @Param request body UpdateCategoryRequest true "Request body to update a category"
// @Router /categories/{category_id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
		return
	}

	var request UpdateCategoryRequest
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

	input := usecases.UpdateCategoryInputDto{
		UserID:     userID,
		CategoryID: request.CategoryID,
		Name:       request.Name,
		Color:      request.Color,
	}

	output, errs := h.categoryFactory.UpdateCategory.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

// DeleteCategory godoc
// @Summary Delete a category
// @Description Delete a category by its ID
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {object} usecases.DeleteCategoryOutputDto
// @Failure 400 {object} util.ProblemDetails
// @Failure 404 {object} util.ProblemDetails
// @Failure 500 {object} util.ProblemDetails
// @Param category_id query string true "Category ID"
// @Router /categories/{category_id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
		return
	}

	categoryID := c.Query("category_id")
	if categoryID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": util.ProblemDetails{
			Type:     "Bad Request",
			Title:    "Missing Category ID",
			Status:   http.StatusBadRequest,
			Detail:   "Category id is required",
			Instance: util.RFC400,
		}})
		return
	}

	input := usecases.DeleteCategoryInputDto{
		UserID:     userID,
		CategoryID: categoryID,
	}

	output, errs := h.categoryFactory.DeleteCategory.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}
