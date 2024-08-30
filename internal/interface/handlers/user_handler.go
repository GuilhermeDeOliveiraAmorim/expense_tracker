package handlers

import (
	"net/http"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/infra/factory"
	usecases "github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/use_cases"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, erros := h.userFactory.CreateUser.Execute(input)
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

func (h *UserHandler) GetUser(c *gin.Context) {
	var input usecases.GetUserInputDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.userFactory.GetUser.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	var input usecases.GetUsersInputDto
	output, err := h.userFactory.GetUsers.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var input usecases.UpdateUserInputDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, erros := h.userFactory.UpdateUser.Execute(input)
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

	c.JSON(http.StatusOK, output)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	var input usecases.DeleteUserInputDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, erros := h.userFactory.DeleteUser.Execute(input)
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

	c.JSON(http.StatusOK, output)
}

func (h *UserHandler) Login(c *gin.Context) {
	var input usecases.LoginInputDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, erros := h.userFactory.Login.Execute(input)
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

	c.JSON(http.StatusOK, output)
}
