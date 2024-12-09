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

// @Summary Create a new user
// @Description Registers a new user in the system
// @Tags Authentication
// @Accept json
// @Produce json
// @Param CreateUserRequest body usecases.CreateUserInputDto true "User data"
// @Success 201 {object} usecases.CreateUserOutputDto
// @Failure 400 {object} util.ProblemDetails "Bad Request"
// @Router /signup [post]
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

// @Summary Get a user by ID
// @Description Retrieves the details of a user by their user_id
// @Tags Users
// @Produce json
// @Param user_id query string true "User ID"
// @Success 200 {object} usecases.GetUserOutputDto
// @Failure 400 {object} util.ProblemDetails "Bad Request"
// @Failure 401 {object} util.ProblemDetails "Unauthorized"
// @Security BearerAuth
// @Router /users [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Query("user_id")
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

// @Summary Get all users
// @Description Retrieves a list of all users
// @Tags Users
// @Produce json
// @Success 200 {array} usecases.GetUsersOutputDto
// @Failure 401 {object} util.ProblemDetails "Unauthorized"
// @Security BearerAuth
// @Router /users/all [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	var input usecases.GetUsersInputDto
	output, errs := h.userFactory.GetUsers.Execute(input)
	if len(errs) > 0 {
		handleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

// @Summary Update a user
// @Description Updates details of a specific user
// @Tags Users
// @Accept json
// @Produce json
// @Param UpdateUserRequest body usecases.UpdateUserInputDto true "User data to update"
// @Success 200 {object} usecases.UpdateUserOutputDto
// @Failure 400 {object} util.ProblemDetails "Bad Request"
// @Failure 401 {object} util.ProblemDetails "Unauthorized"
// @Security BearerAuth
// @Router /users [patch]
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

// @Summary Delete a user by ID
// @Description Deletes a specific user by their user_id
// @Tags Users
// @Produce json
// @Param user_id query string true "User ID"
// @Success 200 {object} usecases.DeleteUserOutputDto
// @Failure 400 {object} util.ProblemDetails "Bad Request"
// @Failure 401 {object} util.ProblemDetails "Unauthorized"
// @Security BearerAuth
// @Router /users [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Query("user_id")
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

// @Summary Login a user
// @Description Authenticates a user and returns a JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param LoginRequest body usecases.LoginInputDto true "User credentials"
// @Success 200 {object} usecases.LoginOutputDto
// @Failure 400 {object} util.ProblemDetails "Bad Request"
// @Router /login [post]
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
