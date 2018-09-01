package main

import (
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"os"
)

func respondWithError(code int, message string, c *gin.Context) {
	resp := map[string]string{"error": message}
	c.JSON(code, resp)
	c.Abort()
}

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// read from header
		token := c.Request.Header.Get("api_token")

		if token == "" {
			respondWithError(401, "API token required", c)
			return
		}
		if token != os.Getenv("API_TOKEN") {
			respondWithError(401, "Invalid API token", c)
			return
		}
		c.Next()
	}
}

func RequestIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		u, _ := uuid.NewV4()
		c.Writer.Header().Set("X-Request-Id", u.String())
		c.Next()
	}
}