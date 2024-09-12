package handlers

import (
	"net/http"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/infra/factory"
	usecases "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/use_cases"
	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	tagFactory *factory.TagFactory
}

func NewTagHandler(factory *factory.TagFactory) *TagHandler {
	return &TagHandler{
		tagFactory: factory,
	}
}

func (h *TagHandler) CreateTag(c *gin.Context) {
	var input usecases.CreateTagInputDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, erros := h.tagFactory.CreateTag.Execute(input)
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

func (h *TagHandler) GetTag(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	tagID := c.Query("tag_id")
	if tagID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tag_id is required"})
		return
	}

	input := usecases.GetTagInputDto{
		UserID: userID,
		TagID:  tagID,
	}

	output, err := h.tagFactory.GetTag.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *TagHandler) GetTags(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	input := usecases.GetTagsInputDto{
		UserID: userID,
	}

	output, errs := h.tagFactory.GetTags.Execute(input)
	if len(errs) > 0 {
		for _, err := range errs {
			if err.Status == 500 {
				c.JSON(err.Status, gin.H{"error": err})
				return
			} else {
				c.JSON(err.Status, gin.H{"error": err})
				return
			}
		}

	}

	c.JSON(http.StatusOK, output)
}

func (h *TagHandler) DeleteTag(c *gin.Context) {
	var input usecases.DeleteTagInputDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.tagFactory.DeleteTag.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, output)
}
