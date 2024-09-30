package middlewares

import (
	"document_microservice/database"
	"document_microservice/helpers"
	"log"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
)

func IsDoctorOrPatient() gin.HandlerFunc {
	return func(c *gin.Context) {
		var idT int
		var isPatientHistory bool
		databaseConn := database.GetConnection()

		uuid, idParam := c.Param("uuid"), c.Param("id")
		if idParam != "" {
			id, err := strconv.Atoi(idParam)
			if err != nil {
				c.JSON(400, gin.H{"message": "parameter id should be a number"})
				c.Abort()
				return
			}
			idT = id
		}

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

		if uuid != "" {
			var uuidInDB string
			row := databaseConn.QueryRow(`SELECT uuid FROM users WHERE username=$1`, claims.Username)

			err = row.Scan(&uuidInDB)
			if err != nil {
				log.Println(err.Error())
			}

			isPatientHistory = uuid == uuidInDB
		}

		if idT != 0 {
			var userUUIDInHistory string
			row := databaseConn.QueryRow(`SELECT pacient_uuid FROM history WHERE id=$1`, idT)

			err = row.Scan(&userUUIDInHistory)
			if err != nil {
				log.Println(err.Error())
			}

			var uuidInDB string
			row2 := databaseConn.QueryRow(`SELECT uuid FROM users WHERE username=$1`, claims.Username)

			err = row2.Scan(&uuidInDB)
			if err != nil {
				log.Println(err.Error())
			}

			isPatientHistory = userUUIDInHistory == uuidInDB
		}

		if !(slices.Contains(claims.Roles, "doctor") || isPatientHistory) {
			c.JSON(401, gin.H{"error": "not allowed"})
			c.Abort()
			return
		}

		c.Set("roles", claims.Roles)
		c.Next()
	}
}
