package handlers

import (
	"net/http"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/infra/factory"
	usecases "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/use_cases"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
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
	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
		return
	}

	var request CreateTagRequest
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

	input := usecases.CreateTagInputDto{
		UserID: userID,
		Name:   request.Name,
		Color:  request.Color,
	}

	output, errs := h.tagFactory.CreateTag.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusCreated, output)
}

func (h *TagHandler) GetTag(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
		return
	}

	tagID := c.Param("tag_id")
	if tagID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": util.ProblemDetails{
			Type:     "Bad Request",
			Title:    "Missing Tag ID",
			Status:   http.StatusBadRequest,
			Detail:   "Tag id is required",
			Instance: util.RFC400,
		}})
		return
	}

	input := usecases.GetTagInputDto{
		UserID: userID,
		TagID:  tagID,
	}

	output, errs := h.tagFactory.GetTag.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *TagHandler) GetTags(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
		return
	}

	input := usecases.GetTagsInputDto{
		UserID: userID,
	}

	output, errs := h.tagFactory.GetTags.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *TagHandler) DeleteTag(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
		return
	}

	tagID := c.Param("tag_id")
	if tagID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": util.ProblemDetails{
			Type:     "Bad Request",
			Title:    "Missing Tag ID",
			Status:   http.StatusBadRequest,
			Detail:   "Tag id is required",
			Instance: util.RFC400,
		}})
		return
	}

	input := usecases.DeleteTagInputDto{
		UserID: userID,
		TagID:  tagID,
	}

	output, errs := h.tagFactory.DeleteTag.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}
