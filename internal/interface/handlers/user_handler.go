package handlers

import (
	"net/http"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/infra/factory"
	usecases "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/use_cases"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userFactory *factory.UserFactory
}

func NewUserHandler(factory *factory.UserFactory) *UserHandler {
	return &UserHandler{
		userFactory: factory,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var input usecases.CreateUserInputDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": util.ProblemDetails{
			Type:     "Bad Request",
			Title:    "Did not bind JSON",
			Status:   http.StatusBadRequest,
			Detail:   err.Error(),
			Instance: util.RFC400,
		}})
		return
	}

	output, errs := h.userFactory.CreateUser.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusCreated, output)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": util.ProblemDetails{
			Type:     "Bad Request",
			Title:    "Missing User ID",
			Status:   http.StatusBadRequest,
			Detail:   "User id is required",
			Instance: util.RFC400,
		}})
		return
	}

	input := usecases.GetUserInputDto{
		UserID: userID,
	}

	output, errs := h.userFactory.GetUser.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	var input usecases.GetUsersInputDto
	output, errs := h.userFactory.GetUsers.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var input usecases.UpdateUserInputDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": util.ProblemDetails{
			Type:     "Bad Request",
			Title:    "Did not bind JSON",
			Status:   http.StatusBadRequest,
			Detail:   err.Error(),
			Instance: util.RFC400,
		}})
		return
	}

	output, errs := h.userFactory.UpdateUser.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": util.ProblemDetails{
			Type:     "Bad Request",
			Title:    "Missing User ID",
			Status:   http.StatusBadRequest,
			Detail:   "User id is required",
			Instance: util.RFC400,
		}})
		return
	}

	input := usecases.DeleteUserInputDto{
		UserID: userID,
	}

	output, errs := h.userFactory.DeleteUser.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *UserHandler) Login(c *gin.Context) {
	var input usecases.LoginInputDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": util.ProblemDetails{
			Type:     "Bad Request",
			Title:    "Did not bind JSON",
			Status:   http.StatusBadRequest,
			Detail:   err.Error(),
			Instance: util.RFC400,
		}})
		return
	}

	output, errs := h.userFactory.Login.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}
