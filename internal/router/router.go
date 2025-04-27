// Package router handles all the HTTP routes for the application
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jinxinyu/go_backend/internal/middleware"
	"gorm.io/gorm"
)

// SetupRouter configures the HTTP router for the application
func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	config := &middleware.CorsOptions{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	r.Use(middleware.NewMiddleware(config))

	// Public routes
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Add more routes here...

	return r
}
