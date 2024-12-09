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

// @Summary Create a new tag
// @Description Create a tag with a name and color
// @Tags Tags
// @Accept json
// @Produce json
// @Param CreateTagRequest body CreateTagRequest true "Tag data"
// @Success 201 {object} usecases.CreateTagOutputDto
// @Failure 400 {object} util.ProblemDetails "Did not bind JSON"
// @Failure 401 {object} util.ProblemDetails "Unauthorized"
// @Security BearerAuth
// @Router /tags [post]
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

// @Summary Get a tag by ID
// @Description Retrieve a tag by its ID
// @Tags Tags
// @Produce json
// @Param tag_id query string true "Tag ID"
// @Success 200 {object} usecases.GetTagOutputDto
// @Failure 400 {object} util.ProblemDetails "Missing Tag ID"
// @Failure 401 {object} util.ProblemDetails "Unauthorized"
// @Security BearerAuth
// @Router /tags/{tag_id} [get]
func (h *TagHandler) GetTag(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
		return
	}

	tagID := c.Query("tag_id")
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

// @Summary Get all tags
// @Description Retrieve all tags for the authenticated user
// @Tags Tags
// @Produce json
// @Success 200 {object} usecases.GetTagsOutputDto
// @Failure 401 {object} util.ProblemDetails "Unauthorized"
// @Security BearerAuth
// @Router /tags [get]
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

// @Summary Delete a tag by ID
// @Description Delete a specific tag by its ID
// @Tags Tags
// @Produce json
// @Param tag_id query string true "Tag ID"
// @Success 200 {object} usecases.DeleteTagOutputDto
// @Failure 400 {object} util.ProblemDetails "Missing Tag ID"
// @Failure 401 {object} util.ProblemDetails "Unauthorized"
// @Security BearerAuth
// @Router /tags/{tag_id} [delete]
func (h *TagHandler) DeleteTag(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
		return
	}

	tagID := c.Query("tag_id")
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

// @Summary Update a tag by ID
// @Description Update a tag's name and color by its ID
// @Tags Tags
// @Accept json
// @Produce json
// @Param UpdateTagRequest body UpdateTagRequest true "Updated tag data"
// @Success 200 {object} usecases.UpdateTagOutputDto
// @Failure 400 {object} util.ProblemDetails "Did not bind JSON"
// @Failure 401 {object} util.ProblemDetails "Unauthorized"
// @Security BearerAuth
// @Router /tags/{tag_id} [put]
func (h *TagHandler) UpdateTag(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(err.Status, gin.H{"error": err})
		return
	}

	var request UpdateTagRequest
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

	input := usecases.UpdateTagInputDto{
		UserID: userID,
		TagID:  request.TagID,
		Name:   request.Name,
		Color:  request.Color,
	}

	output, errs := h.tagFactory.UpdateTag.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}
