package middlewares

import (
	"account_microservice/helpers"
	"github.com/gin-gonic/gin"
	"slices"
)

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("tokenAccess")
		if err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		claims, err := helpers.ParseToken(cookie)
		if err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if !slices.Contains(claims.Roles, "admin") {
			c.JSON(401, gin.H{"error": "You are not an admin!"})
			c.Abort()
			return
		}

		c.Set("roles", claims.Roles)
		c.Next()
	}
}