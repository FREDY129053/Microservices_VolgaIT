package controllers

import (
	"hospital_microservice/database"
	"hospital_microservice/models"
	"hospital_microservice/helpers"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

var databaseConn = database.GetConnection()

func GetAllHospitals(c *gin.Context) {
	// Получение параметров
	fromParam := c.Query("from")
	countParam := c.Query("count")

	// Проверка параметров
	from, err := strconv.Atoi(fromParam)
	if err != nil {
		c.JSON(400, gin.H{"message": "Parameter from should be a number"})
		c.Abort()
		return
	}
	count, err := strconv.Atoi(countParam)
	if err != nil {
		c.JSON(400, gin.H{"message": "Parameter count should be a number"})
		c.Abort()
		return
	}

	var hospitals []models.HospitalInfo

	rows, err := databaseConn.Query("SELECT * FROM hospital LIMIT $1 OFFSET $2", count, from-1)
	if err != nil {
		c.JSON(501, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	for rows.Next() {
		hospital := models.HospitalInfo{}
		err := rows.Scan(&hospital.UUID, &hospital.Name, &hospital.Address, &hospital.ContactPhone)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		// Получение комнат больницы
		var rooms []string
		rows, err := databaseConn.Query("SELECT room FROM hospital_rooms WHERE hospital_uuid=$1", hospital.UUID)
		if err != nil {
			c.JSON(501, gin.H{"message": err.Error()})
			c.Abort()
			return
		}
		// Добавляем комнаты в массив
		for rows.Next() {
			var room string
			err := rows.Scan(&room)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			rooms = append(rooms, room)
		}

		hospital.Rooms = rooms
		hospitals = append(hospitals, hospital)
	}

	if len(hospitals) == 0 {
		c.JSON(200, gin.H{"message": "No hospital in database"})
		c.Abort()
		return
	}

	c.JSON(200, hospitals)
}

func GetHospitalInfo(c *gin.Context) {
	var hospitalInfo models.HospitalInfo
	hospitalUUID := c.Param("uuid")

	row := databaseConn.QueryRow("SELECT * FROM hospital WHERE uuid=$1", hospitalUUID)
	if err := row.Scan(&hospitalInfo.UUID, &hospitalInfo.Name, &hospitalInfo.Address, &hospitalInfo.ContactPhone); err != nil {
		c.JSON(400, gin.H{"message": "Cannot find hospital"})
		c.Abort()
		return
	}

	var rooms []string
	roomsRows, err := databaseConn.Query("SELECT room FROM hospital_rooms WHERE hospital_uuid=$1", hospitalInfo.UUID)
	if err != nil {
		panic(err)
	}
	for roomsRows.Next() {
		var room string
		err := roomsRows.Scan(&room)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		rooms = append(rooms, room)
	}

	hospitalInfo.Rooms = rooms

	c.JSON(200, hospitalInfo)
}

func GetHospitalRooms(c *gin.Context) {
	var rooms []string
	hospitalUUID := c.Param("uuid")

	roomsRows, err := databaseConn.Query("SELECT room FROM hospital_rooms WHERE hospital_uuid=$1", hospitalUUID)
	if err != nil {
		c.JSON(404, gin.H{"message": "Hospital not found"})
	}
	for roomsRows.Next() {
		var room string
		err := roomsRows.Scan(&room)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		rooms = append(rooms, room)
	}

	c.JSON(200, rooms)
}

func AddHospital(c *gin.Context) {
	var hospitalInfo models.AddHospitalInfo

	if err := c.ShouldBindJSON(&hospitalInfo); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	// Проверка на наличие больницы в БД
	var existingHospital models.AddHospitalInfo
	row := databaseConn.QueryRow("SELECT name,address,contact_phone FROM hospital WHERE name=$1 OR address=$2 OR contact_phone=$3", hospitalInfo.Name, hospitalInfo.Address, hospitalInfo.ContactPhone)
	if err := row.Scan(&existingHospital.Name, &existingHospital.Address, &existingHospital.ContactPhone); err == nil {
		c.JSON(400, gin.H{"message": "Hospital already exists"})
		c.Abort()
		return
	}

	hospitalUUID := helpers.GenerateUUID()
	_, err := databaseConn.Exec("INSERT INTO hospital (uuid, name, address, contact_phone) VALUES($1, $2, $3, $4)", hospitalUUID, hospitalInfo.Name, hospitalInfo.Address, hospitalInfo.ContactPhone)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		c.Abort()
		return
	}
	for _, room := range hospitalInfo.Rooms {
		_, err := databaseConn.Exec("INSERT INTO hospital_rooms (hospital_uuid, room) VALUES ($1, $2)", hospitalUUID, room)
		if err != nil {
			c.JSON(501, gin.H{"message": err.Error()})
			c.Abort()
			return
		}
	}

	c.JSON(200, gin.H{"message": "Hospital created successfully"})
}

func UpdateHospital(c *gin.Context) {
	var hospitalInfo models.AddHospitalInfo
	hospitalUUID := c.Param("uuid")

	if err := c.ShouldBindJSON(&hospitalInfo); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		c.Abort()
		return
	} 

	// Для начала удалим все комнаты у больницы в БД
	_, err := databaseConn.Exec("DELETE FROM hospital_rooms WHERE hospital_uuid=$1", hospitalUUID)
	if err != nil {
		c.JSON(404, gin.H{"message": "Hospital not found"})
		c.Abort()
		return
	}

	// Вставка комнат и изменение данных
	_, err2 := databaseConn.Exec(`
		UPDATE hospital 
		SET name=$1, address=$2, contact_phone=$3
		WHERE uuid=$4`,
		hospitalInfo.Name, hospitalInfo.Address, hospitalInfo.ContactPhone, hospitalUUID,
	)
	if err2 != nil {
		c.JSON(500, gin.H{"message": err2.Error()})
		c.Abort()
		return
	}

	for _, room := range hospitalInfo.Rooms {
		_, err := databaseConn.Exec(`
			INSERT INTO hospital_rooms (hospital_uuid, room) 
			VALUES ($1, $2)`, 
			hospitalUUID, room,
		)
		if err != nil {
			c.JSON(501, gin.H{"message": err.Error()})
			c.Abort()
			return
		}
	}

	c.JSON(200, gin.H{"message": "Hospital updated successfully"})
}

func DeleteHospital(c *gin.Context) {
	hospitalUUID := c.Param("uuid")
	_, err := databaseConn.Exec("DELETE FROM hospital WHERE uuid=$1", hospitalUUID)
	if err != nil {
		c.JSON(404, gin.H{"message": "Oi, hospital not found"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"message": "Hospital deleted successfully"})
}