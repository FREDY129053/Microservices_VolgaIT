package middlewares

import (
	"timetable_microservice/helpers"
	"timetable_microservice/database"
	"github.com/gin-gonic/gin"
	"slices"
	"log"
	"strconv"
)

func IsAdminOrManagerOrPatient() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("tokenAccess")
		if err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		appointmentIDParam := c.Param("id")
		// Проверка параметров
		appointmentID, err := strconv.Atoi(appointmentIDParam)
		if err != nil {
			c.JSON(400, gin.H{"message": "parameter id should be a number"})
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

		// Получение username записи из БД
		databaseConn := database.GetConnection()
		var usernameInDB string
		row := databaseConn.QueryRow("SELECT pacient_username FROM appointments WHERE id=$1", appointmentID)
		if err := row.Scan(&usernameInDB); err != nil {
			c.JSON(400, gin.H{"message": "Cannot find appointment note"})
			c.Abort()
			return
		}
		

		if !(slices.Contains(claims.Roles, "admin") ||
				slices.Contains(claims.Roles, "manager") ||
				(usernameInDB == claims.Username)) {
			c.JSON(401, gin.H{"error": "not allowed"})
			c.Abort()
			return
		}

		c.Set("roles", claims.Roles)
		c.Next()
	}
}