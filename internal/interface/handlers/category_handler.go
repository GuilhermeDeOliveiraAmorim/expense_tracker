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

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	userID, err := getUserID(c)
	if len(err) > 0 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
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

	output, erros := h.categoryFactory.CreateCategory.Execute(input)
	if len(erros) > 0 {
		for _, err := range erros {
			if err.Status == 500 {
				c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
				return
			} else {
				c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
				return
			}
		}

	}

	c.JSON(http.StatusCreated, output)
}

func (h *CategoryHandler) GetCategory(c *gin.Context) {
	userID, err := getUserID(c)
	if len(err) > 0 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	categoryID := c.Query("category_id")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category_id is required"})
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

	output, erros := h.categoryFactory.GetCategory.Execute(input)
	if len(erros) > 0 {
		for _, err := range erros {
			if err.Status == 500 {
				c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
				return
			} else {
				c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
				return
			}
		}
	}

	c.JSON(http.StatusOK, output)
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	userID, err := getUserID(c)
	if len(err) > 0 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	input := usecases.GetCategoriesInputDto{
		UserID: userID,
	}

	output, errs := h.categoryFactory.GetCategories.Execute(input)
	if len(errs) > 0 {
		for _, err := range errs {
			if err.Status == 500 {
				c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
				return
			} else {
				c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
				return
			}
		}

	}

	c.JSON(http.StatusOK, output)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	userID, err := getUserID(c)
	if len(err) > 0 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
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
		for _, err := range errs {
			if err.Status == 500 {
				c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
				return
			} else {
				c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
				return
			}
		}

	}

	c.JSON(http.StatusOK, output)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	userID, err := getUserID(c)
	if len(err) > 0 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	categoryID := c.Query("category_id")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category_id is required"})
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

	output, erros := h.categoryFactory.DeleteCategory.Execute(input)
	if len(erros) > 0 {
		for _, err := range erros {
			if err.Status == 500 {
				c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
				return
			} else {
				c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
				return
			}
		}
	}

	c.JSON(http.StatusOK, output)
}
