// Package router handles all the HTTP routes for the application
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jinxinyu/go_backend/internal/auth"
	"github.com/jinxinyu/go_backend/internal/middleware"
)

// SetupRouter configures the HTTP router for the application
func SetupRouter(authService *auth.Service) *gin.Engine {
	r := gin.Default()
	config := &middleware.CorsOptions{
		AllowAllOrigins:  []string{"http://localhost:3000"},
		AllowAllMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowAllHeaders:  []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"},
		AllowCredentials: true,
	}
	r.Use(middleware.NewMiddleware(config))

	// Public routes
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Register auth routes
	apiv1 := r.Group("/api/v1")
	auth.RegisterUserRoutes(apiv1, authService)
	// Add more routes here...

	return r
}
