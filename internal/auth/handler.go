package auth

import (
	"context"
	"log"
	"net/http"
	"time"

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
	// 创建一个更长超时的上下文(60秒)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
	defer cancel()

	var req api.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("开始创建用户: %s", req.Email)
	startTime := time.Now()

	user, err := h.service.RegisterUser(ctx, &req)

	log.Printf("创建用户操作耗时: %v", time.Since(startTime))
	if err != nil {
		log.Printf("创建用户失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

func (h *Handler) LoginUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	var req api.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.service.LoginUser(ctx, &req)
	if err != nil {
		log.Printf("登录失败: %v", err)

		// 提供更友好的错误消息
		if err.Error() == "invalid credentials" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码不正确"})
			return
		}

		// 检查是否为特定的数据库错误
		if err.Error() == "failed to get user by email: record not found" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码不正确"})
			return
		}

		// 其他服务器内部错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": "登录时发生错误，请稍后再试"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}
