package handler

import (
	"net/http"
	"strconv"

	"github.com/EngenMe/go-clean-arch-auth/internal/core/features/mediator"
	userCommands "github.com/EngenMe/go-clean-arch-auth/internal/core/features/user/command/requests"
	userQueries "github.com/EngenMe/go-clean-arch-auth/internal/core/features/user/query/requests"
	userResponse "github.com/EngenMe/go-clean-arch-auth/internal/core/features/user/query/response"
	"github.com/EngenMe/go-clean-arch-auth/internal/core/http/middleware"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	mediator *mediator.Mediator
}

func NewUserHandler(mediator *mediator.Mediator) *UserHandler {
	return &UserHandler{mediator: mediator}
}

// TODO: inject handle errors in 3rd provider
func (h *UserHandler) Register(c *gin.Context) {
	var req userCommands.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := h.mediator.Send(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//TODO: inject 3rd party to handle response
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req userQueries.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.mediator.Send(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	user, ok := result.(*userResponse.UserResponse)
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "invalid response type"},
		)
		return
	}

	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to generate token"},
		)
		return
	}

	c.JSON(
		http.StatusOK, gin.H{
			"message": "Login successful",
			"token":   token,
			"user":    user,
		},
	)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	req := &userQueries.FindUserByIdRequest{ID: uint(id)}
	result, err := h.mediator.Send(req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	user, ok := result.(*userResponse.UserResponse)
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "invalid response type"},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists || userID.(uint) != uint(id) {
		c.JSON(
			http.StatusForbidden,
			gin.H{"error": "you can only update your own account"},
		)
		return
	}

	var userUpdate userCommands.UpdateUserRequest
	if err := c.ShouldBindJSON(&userUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userUpdate.ID = uint(id)
	if _, err := h.mediator.Send(&userUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch updated user
	req := &userQueries.FindUserByIdRequest{ID: uint(id)}
	result, err := h.mediator.Send(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedUser, ok := result.(*userResponse.UserResponse)
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "invalid response type"},
		)
		return
	}

	c.JSON(
		http.StatusOK, gin.H{
			"message": "User updated successfully",
			"user":    updatedUser,
		},
	)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists || userID.(uint) != uint(id) {
		c.JSON(
			http.StatusForbidden,
			gin.H{"error": "you can only delete your own account"},
		)
		return
	}

	req := &userCommands.DeleteUserRequest{ID: uint(id)}
	if _, err := h.mediator.Send(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	req := &userQueries.ListUsersRequest{Page: page, Limit: limit}
	result, err := h.mediator.Send(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	users, ok := result.([]*userResponse.UserResponse)
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "invalid response type"},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}
