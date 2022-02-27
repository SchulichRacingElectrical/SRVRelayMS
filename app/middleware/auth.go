package middleware

import (
	"database-ms/config"

	"github.com/gin-gonic/gin"
)

// func success(c *gin.Context, organization *map[string]interface{}) {
// 	c.Set("organization", *organization)
// 	c.Next()
// }

func respondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{"error": message})
}

// Returns the organization associated with the authorization
func AuthorizationMiddleware(conf *config.Configuration) gin.HandlerFunc {

	adminToken := conf.AdminKey

	return func(c *gin.Context) {

		token := c.Request.Header.Get("api_token")

		if token == "" {
			respondWithError(c, 401, "API token required")
			return
		}

		if token == adminToken {
			c.Next()
		}
	}
}
