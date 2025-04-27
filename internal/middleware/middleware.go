package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type CorsOptions struct {
	AllowAllOrigins  []string
	AllowAllMethods  []string
	AllowAllHeaders  []string
	AllowCredentials bool
}

func NewMiddleware(options *CorsOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		methods := "GET, POST, PUT, DELETE, OPTIONS"
		// Allow all methods if specified
		if len(options.AllowAllMethods) > 0 {
			methods = strings.Join(options.AllowAllMethods, ",")
		}
		headers := "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"
		// Allow all headers if specified
		if len(options.AllowAllHeaders) > 0 {
			headers = strings.Join(options.AllowAllHeaders, ",")
		}

		reqOrigin := c.Request.Header.Get("Origin")
		allowed := false
		allowedValues := ""
		// Check if the origin is allowed
		if len(options.AllowAllOrigins) == 0 {
			allowed = true
			allowedValues = "*"
		} else if len(options.AllowAllOrigins) == 1 && options.AllowAllOrigins[0] == "*" {
			allowed = true
			allowedValues = "*"
		} else {
			for _, origin := range options.AllowAllOrigins {
				if origin == reqOrigin {
					allowed = true
					allowedValues = reqOrigin
					break
				}
			}
		}
		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowedValues)
			c.Writer.Header().Set("Access-Control-Allow-Methods", methods)
			c.Writer.Header().Set("Access-Control-Allow-Headers", headers)
			if options.AllowCredentials {
				if allowedValues == "*" && reqOrigin != "" {
					c.Writer.Header().Set("Access-Control-Allow-Origin", reqOrigin)
				}
			}
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		} else {

		}
		if c.Request.Method == "OPTIONS" {
			c.Writer.Header().Set("Access-Control-Allow-Methods", methods)
			c.Writer.Header().Set("Access-Control-Allow-Headers", headers)
			if options.AllowCredentials {
				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			}
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}

}
