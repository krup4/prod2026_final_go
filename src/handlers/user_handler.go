package handlers

import (
	"errors"
	"net/http"

	"backend/exceptions"
	"backend/handlers/requests"
	"backend/services"
	"backend/services/inputs"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Create(c *gin.Context) {
	var req requests.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json payload"})
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), inputs.CreateUserInput{
		Login: req.Login,
		Role:  req.Role,
	})
	if err != nil {
		if errors.Is(err, exceptions.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "login and role are required"})
			return
		}
		if errors.Is(err, exceptions.ErrUserExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) List(c *gin.Context) {
	users, err := h.userService.GetUsersList(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get list of users"})
	}

	c.JSON(http.StatusOK, users)
}
