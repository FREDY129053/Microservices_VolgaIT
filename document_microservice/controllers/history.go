package controllers

import (
	"document_microservice/database"
	"document_microservice/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var databaseConn = database.GetConnection()

// Отправка запроса на микросервис больниц
func _IsHospitalExist(uuid, token string) bool {
	url := fmt.Sprintf("http://0.0.0.0:8082/api/Hospitals/%s", uuid)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}

	// Задаем куки с токеном
	req.AddCookie(&http.Cookie{
		Name: "tokenAccess",
		Value: token,
	})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}

	return resp.StatusCode == 200
}

// Отправка запроса на микросервис аккаунтов для проверки доктора
func _IsDoctorExist(uuid, token string) bool {
	url := fmt.Sprintf("http://0.0.0.0:8081/api/Accounts/Doctors/%s", uuid)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}

	// Задаем куки с токеном
	req.AddCookie(&http.Cookie{
		Name: "tokenAccess",
		Value: token,
	})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}

	return resp.StatusCode == 200
}

// Проверка на сущестование комнаты
func _IsRoomExist(hospitalUUID, token, room string) bool {
	url := fmt.Sprintf("http://0.0.0.0:8082/api/Hospitals/%s/Rooms", hospitalUUID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}

	// Задаем куки с токеном
	req.AddCookie(&http.Cookie{
		Name: "tokenAccess",
		Value: token,
	})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	return strings.Contains(string(body), room)
}

// Проверка что patientID это акк с ролью User
func _IsPatientUser(patientUUID string) bool {
	var idTemp int
	row := databaseConn.QueryRow("SELECT id FROM user_and_roles WHERE user_uuid=$1 AND role=$2", patientUUID, "user")
	if err := row.Scan(&idTemp); err != nil {
		return false
	}

	return true
}


// GetAllAccountHistories godoc
// GetAllAccountHistories получение всех посещений и назначений аккаунта
// @Summary Получение всех посещений и назначений аккаунта
// @Description Получение всех посещений и назначений аккаунта по ID. Только врачи и тот, кому принадлежит история
// @Tags Documents
// @Accept json
// @Produce json
// @Param uuid path string true "UUID аккаунта(пациента)"
// @Success 200 {object} map[string][]models.FullHistory "Все посещения и назначения аккаунта"
// @Failure 400 {object} map[string]string "Invalid request"
// @Router /Account/{uuid} [get]
// @Security ApiKeyAuth
func GetAllAccountHistories(c *gin.Context) {
	var allHistories []models.FullHistory
	uuidPatient := c.Param("uuid")

	rows, err := databaseConn.Query(`SELECT * FROM history WHERE pacient_uuid=$1`, uuidPatient)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid request"})
		c.Abort()
		return
	}

	for rows.Next() {
		history := models.FullHistory{}
		// Считали данные	
		err := rows.Scan(&history.ID, &history.Date, &history.PatientUUID, &history.HospitalUUID, &history.DoctorUUID, &history.Room, &history.Data)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		allHistories = append(allHistories, history)
	}

	c.JSON(200, gin.H{"histories": allHistories})
}


// GetHistory godoc
// GetHistory получение конкретной истории посещения и назначений
// @Summary Получение конкретной истории посещения и назначений
// @Description Получение конкретной истории посещения и назначений по ID. Только врачи и тот, кому принадлежит история
// @Tags Documents
// @Accept json
// @Produce json
// @Param id path string true "ID истории"
// @Success 200 {object} map[string][]models.FullHistory "Информация о конкретной истории"
// @Failure 400 {object} map[string]string "Invalid request"
// @Router /{id} [get]
// @Security ApiKeyAuth
func GetHistory(c *gin.Context) {
	var history models.FullHistory
	idParam := c.Param("id")

	// Проверка параметров
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"message": "parameter id should be a number"})
		c.Abort()
		return
	}

	row := databaseConn.QueryRow(`SELECT * FROM history WHERE id=$1`, id)

	err = row.Scan(&history.ID, &history.Date, &history.PatientUUID, &history.HospitalUUID, &history.DoctorUUID, &history.Room, &history.Data)
	if err != nil {
		log.Println(err.Error())
	}

	c.JSON(200, gin.H{"history": history})
}


// AddNewHistory godoc
// AddNewHistory добавление истории посещения
// @Summary Добавление новой истории посещений
// @Description Добавление новой истории посещений. Только админы, менеджеры и врачи. PatientUUID - аккаунт с ролью user
// @Tags Documents
// @Accept json
// @Produce json
// @Param history body models.HistoryInfo true "Информация об истории"
// @Success 200 {object} map[string]string "message": "history added successfully"
// @Failure 400 {object} map[string]string "message": "invalid request/patient must be user"
// @Failure 404 {object} map[string]string "message": "hospital/doctor/room not found"
// @Failure 500 {object} map[string]string "message": "internal server error"
// @Router / [post]
// @Security ApiKeyAuth
func AddNewHistory(c *gin.Context) {
	var newHistory models.HistoryInfo
	token, err := c.Cookie("tokenAccess")
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&newHistory); err != nil {
		c.JSON(401, gin.H{"message": err.Error()})
		return
	}

	if !_IsHospitalExist(newHistory.HospitalUUID, token) {
		c.JSON(404, gin.H{"message": "hospital not found"})
		return
	}

	if !_IsDoctorExist(newHistory.DoctorUUID, token) {
		c.JSON(404, gin.H{"message": "doctor not found"})
		return
	}

	if !_IsRoomExist(newHistory.HospitalUUID, token, newHistory.Room) {
		c.JSON(404, gin.H{"message": "room not found"})
		return
	}

	if !_IsPatientUser(newHistory.PatientUUID) {
		c.JSON(400, gin.H{"message": "patient must be user"})
		return
	}

	_, err = databaseConn.Exec(`
		INSERT INTO history (date, pacient_uuid, hospital_uuid, doctor_uuid, room, data)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		newHistory.Date, newHistory.PatientUUID, newHistory.HospitalUUID, newHistory.DoctorUUID, newHistory.Room, newHistory.Data,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "history added successfully"})
}


// UpdateHistory godoc
// UpdateHistory обновляет историю посещения по id
// @Summary Обновление истории посещения по ID
// @Description Обновление конкретной истории посещения по ID. Только админы, менеджеры и врачи.
// @Tags Documents
// @Accept json
// @Produce json
// @Param id path string true "ID истории"
// @Param history body models.HistoryInfo true "Информация ою истории"
// @Success 200 {object} map[string]string "message": "history updated successfully"
// @Failure 400 {object} map[string]string "message": "invalid request/parameter id should be a number"
// @Failure 404 {object} map[string]string "message": "hospital/doctor/room not found"
// @Failure 500 {object} map[string]string "message": "internal server error"
// @Router /{id} [put]
// @Security ApiKeyAuth
func UpdateHistory(c *gin.Context) {
	var historyInfo models.HistoryInfo
	idParam := c.Param("id")

	// Проверка параметров
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"message": "parameter id should be a number"})
		c.Abort()
		return
	}

	token, err := c.Cookie("tokenAccess")
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&historyInfo); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	
	if !_IsHospitalExist(historyInfo.HospitalUUID, token) {
		c.JSON(404, gin.H{"message": "hospital not found"})
		return
	}

	if !_IsDoctorExist(historyInfo.DoctorUUID, token) {
		c.JSON(404, gin.H{"message": "doctor not found"})
		return
	}

	if !_IsRoomExist(historyInfo.HospitalUUID, token, historyInfo.Room) {
		c.JSON(404, gin.H{"message": "room not found"})
		return
	}

	if !_IsPatientUser(historyInfo.PatientUUID) {
		c.JSON(400, gin.H{"message": "patient must be user"})
		return
	}

	_, err = databaseConn.Exec(`
		UPDATE history 
		SET date=$1,
		pacient_uuid=$2, 
		hospital_uuid=$3,
		doctor_uuid=$4, 
		room=$5,
		data=$6,
		WHERE id=$7`,
		historyInfo.Date, historyInfo.PatientUUID, historyInfo.HospitalUUID, historyInfo.DoctorUUID, historyInfo.Room, historyInfo.Data, id,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "history updated successfully"})
}