package auth

import (
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.RouterGroup, service *Service) {
	handler := NewHandler(service)

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", handler.RegisterUser)
		authRoutes.POST("/login", handler.LoginUser)
	}
}
