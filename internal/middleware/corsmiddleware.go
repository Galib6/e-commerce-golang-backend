package middleware

import "github.com/gin-gonic/gin"

func CORSMiddleware() gin.HandlerFunc {

	allowedOrigins := map[string]bool{
		"http://localhost:3000":         true, // dev
		"https://your-frontend.com":     true, // prod
		"https://www.your-frontend.com": true,
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Vary", "Origin")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true") // if you use cookies
		}
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, Authorization, Origin, Accept")
		c.Writer.Header().Set("Access-Control-Allow-Methods",
			"GET, POST, PUT, PATCH, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
