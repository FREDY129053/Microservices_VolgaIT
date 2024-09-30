package middlewares

import (
	"timetable_microservice/helpers"
	"github.com/gin-gonic/gin"
	"slices"
	"log"
)

func IsAdminOrManagerOrDoctor() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("tokenAccess")
		if err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		claims, err := helpers.ParseToken(cookie)
		if err != nil {
			log.Println(err.Error())
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		

		if !(slices.Contains(claims.Roles, "admin") ||
				slices.Contains(claims.Roles, "manager") ||
				slices.Contains(claims.Roles, "doctor")) {
			c.JSON(401, gin.H{"error": "not allowed"})
			c.Abort()
			return
		}

		c.Set("roles", claims.Roles)
		c.Next()
	}
}