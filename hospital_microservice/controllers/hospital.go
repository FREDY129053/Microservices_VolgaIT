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

// GetAllHospitals godoc
// GetAllHospitals получение больниц в базе данных
// @Summary Получение больниц в базе данных
// @Description Получение определенного числа больниц в базе данных. Только авторизованные пользователи.
// @Tags Hospitals
// @Accept json
// @Produce json
// @Param from path string true "Начало выборки(порядковый номер)"
// @Param count path string true "Размер выборки"
// @Success 200 {object} map[string][]models.HospitalInfo "Все больницы"
// @Failure 400 {object} map[string]string "Parameter from/count should be a number"
// @Failure 501 {object} map[string]string "Internal Server Error"
// @Router / [get]
// @Security ApiKeyAuth
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


// GetHospitalInfo godoc
// GetHospitalInfo получение информации о больнице
// @Summary Получение информации о больнице
// @Description Получение информации о больнице по UUID. Только авторизованные пользователи
// @Tags Hospitals
// @Accept json
// @Produce json
// @Param uuid path string true "UUD больницы"
// @Success 200 {object} []models.HospitalInfo "Информация о больнице"
// @Failure 404 {object} map[string]string "Cannot find hospital"
// @Router /{uuid} [get]
// @Security ApiKeyAuth
func GetHospitalInfo(c *gin.Context) {
	var hospitalInfo models.HospitalInfo
	hospitalUUID := c.Param("uuid")

	row := databaseConn.QueryRow("SELECT * FROM hospital WHERE uuid=$1", hospitalUUID)
	if err := row.Scan(&hospitalInfo.UUID, &hospitalInfo.Name, &hospitalInfo.Address, &hospitalInfo.ContactPhone); err != nil {
		c.JSON(404, gin.H{"message": "Cannot find hospital"})
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


// GetHospitalRooms godoc
// GetHospitalRooms получение кабинетов больницы
// @Summary Получение списка кабинетов больницы
// @Description Получение списка всех кабинетов больницы по ее UUID. Только авторизованные пользователи
// @Tags Hospitals
// @Accept json
// @Produce json
// @Param uuid path string true "UUD больницы"
// @Success 200 {object} []string "Список кабинетов больницы"
// @Failure 404 {object} map[string]string "Hospital not found"
// @Router /{uuid}/Rooms [get]
// @Security ApiKeyAuth
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


// AddHospital godoc
// AddHospital добавление больницы
// @Summary Добавление больницы в базу данных
// @Description Добавление больницы с переданной инофрмацией в базу данных. Только админы
// @Tags Hospitals
// @Accept json
// @Produce json
// @Param hospital body models.AddHospitalInfo true "Информация о больнице"
// @Success 200 {object} map[string]string "Hospital created successfully"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal Sever Error"
// @Failure 501 {object} map[string]string "Internal Sever Error"
// @Router / [post]
// @Security ApiKeyAuth
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


// UpdateHospital godoc
// UpdateHospital изменение больницы
// @Summary Изменение больницы в базе данных
// @Description Изменение больницы с переданной инофрмацией в базе данных. Только админы
// @Tags Hospitals
// @Accept json
// @Produce json
// @Param uuid path string true "UUD больницы"
// @Param hospital body models.AddHospitalInfo true "Информация о больнице"
// @Success 200 {object} map[string]string "Hospital updated successfully"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "Hospital not found"
// @Failure 500 {object} map[string]string "Internal Sever Error"
// @Failure 501 {object} map[string]string "Internal Sever Error"
// @Router /{uuid} [put]
// @Security ApiKeyAuth
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


// DeleteHospital godoc
// DeleteHospital удалить больницу
// @Summary Удаление больницы
// @Description Удаление записи о больницу по UUID больницы. Только админы
// @Tags Hospitals
// @Accept json
// @Produce json
// @Param uuid path string true "UUID больницы"
// @Success 200 {object} map[string]string "Hospital deleted successfully"
// @Failure 404 {object} map[string]string "Oi, hospital not found"
// @Router /{uuid} [delete]
// @Security ApiKeyAuth
func DeleteHospital(c *gin.Context) {
	hospitalUUID := c.Param("uuid")

	var dbInfo string
	row := databaseConn.QueryRow("SELECT name FROM hospital WHERE uuid=$1", hospitalUUID)
	if err := row.Scan(&dbInfo); err != nil {
		c.JSON(404, gin.H{"message": "Oi, hospital not found"})
		c.Abort()
		return
	}

	_, err := databaseConn.Exec("DELETE FROM hospital WHERE uuid=$1", hospitalUUID)
	if err != nil {
		c.JSON(404, gin.H{"message": "Oi, hospital not found"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"message": "Hospital deleted successfully"})
}