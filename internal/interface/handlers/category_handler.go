package handlers

import (
	"net/http"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/infra/factory"
	usecases "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/use_cases"
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
	var input usecases.CreateCategoryInputDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, erros := h.categoryFactory.CreateCategory.Execute(input)
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

func (h *CategoryHandler) GetCategory(c *gin.Context) {
	var input usecases.GetCategoryInputDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.categoryFactory.GetCategory.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	var input usecases.GetCategoriesInputDto
	output, err := h.categoryFactory.GetCategories.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	var input usecases.UpdateCategoryInputDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.categoryFactory.UpdateCategory.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	var input usecases.DeleteCategoryInputDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.categoryFactory.DeleteCategory.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, output)
}
