package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinxinyu/go_backend/internal/api"
)

type Handler struct {
	service *Service
}

// NewHandler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterUser(c *gin.Context) {
	ctx := c.Request.Context()

	var req api.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.RegisterUser(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

func (h *Handler) LoginUser(c *gin.Context) {
	ctx := c.Request.Context()

	var req api.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.service.LoginUser(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}
