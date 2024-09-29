package controllers

import (
	"time"
	"fmt"
	"net/http"
	"io"
	"log"
	"strconv"
	"strings"
	"timetable_microservice/database"
	"timetable_microservice/models"

	"github.com/gin-gonic/gin"
)

var databaseConn = database.GetConnection()
const iso8601Layout = "2006-01-02T15:04:05Z07:00"

// Проверка времени в Get-ах
func _IsParamsValid(from, to string) bool {
	_, err := time.Parse(iso8601Layout, from)
	if err != nil {
		return false
	}

	_, err2 := time.Parse(iso8601Layout, to)
	
	return err2 == nil 
}

// Проверки временных промежутков
func _TimeValidate(from, to time.Time) error {
	if from.Minute() % 30 != 0 || from.Second() != 0 {
		return fmt.Errorf("invalid format for date 'from'")
	}
	if to.Minute() % 30 != 0 || to.Second() != 0 {
		return fmt.Errorf("invalid format for date 'to'")
	}
	if to.Before(from) {
		return fmt.Errorf("'to' must be greater than 'from'")
	}
	if to.Sub(from).Hours() > 12 {
		return fmt.Errorf("the difference between 'to' and 'from' must not exceed 12 hours")
	}

	return nil
}

// Отправка запроса на микросервис больниц
func _IsHospitalExist(uuid, token string) bool {
	url := fmt.Sprintf("http://localhost:8082/api/Hospitals/%s", uuid)
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
	url := fmt.Sprintf("http://localhost:8081/api/Accounts/Doctors/%s", uuid)
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
	url := fmt.Sprintf("http://localhost:8082/api/Hospitals/%s/Rooms", hospitalUUID)
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


func AddNewNote(c *gin.Context) {
	var noteInfo models.Timetable

	token, err := c.Cookie("tokenAccess")
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&noteInfo); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	if err := _TimeValidate(noteInfo.From, noteInfo.To); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	
	if !_IsHospitalExist(noteInfo.HospitalUUID, token) {
		c.JSON(404, gin.H{"message": "hospital not found"})
		return
	}

	if !_IsDoctorExist(noteInfo.DoctorUUID, token) {
		c.JSON(404, gin.H{"message": "doctor not found"})
		return
	}

	if !_IsRoomExist(noteInfo.HospitalUUID, token, noteInfo.Room) {
		c.JSON(404, gin.H{"message": "room not found"})
		return
	}

	_, err = databaseConn.Exec(`
		INSERT INTO timetable (hospital_uuid, doctor_uuid, time_from, time_to, room)
		VALUES ($1, $2, $3, $4, $5)`,
		noteInfo.HospitalUUID, noteInfo.DoctorUUID, noteInfo.From, noteInfo.To, noteInfo.Room,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "note added successfully"})
}

func UpdateNote(c *gin.Context) {
	var noteInfo models.Timetable
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

	if err := c.ShouldBindJSON(&noteInfo); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	if err := _TimeValidate(noteInfo.From, noteInfo.To); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	
	if !_IsHospitalExist(noteInfo.HospitalUUID, token) {
		c.JSON(404, gin.H{"message": "hospital not found"})
		return
	}

	if !_IsDoctorExist(noteInfo.DoctorUUID, token) {
		c.JSON(404, gin.H{"message": "doctor not found"})
		return
	}

	if !_IsRoomExist(noteInfo.HospitalUUID, token, noteInfo.Room) {
		c.JSON(404, gin.H{"message": "room not found"})
		return
	}

	// Проверка есть ли записавшийся на прием
	var appointmentExist int
	row := databaseConn.QueryRow("SELECT id FROM appointments WHERE timetable_id=$1", id)
	if err := row.Scan(&appointmentExist); err == nil {
		c.JSON(400, gin.H{"message": "appointments not empty"})
		return
	}

	_, err = databaseConn.Exec(`
		UPDATE timetable 
		SET hospital_uuid=$1,
		doctor_uuid=$2, 
		time_from=$3,
		time_to=$4, 
		room=$5
		WHERE id=$6`,
		noteInfo.HospitalUUID, noteInfo.DoctorUUID, noteInfo.From, noteInfo.To, noteInfo.Room, id,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "note updated successfully"})
}

func DeleteByID(c *gin.Context) {
	idParam := c.Param("id")

	// Проверка параметров
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"message": "parameter id should be a number"})
		c.Abort()
		return
	}

	_, err = databaseConn.Exec("DELETE FROM timetable WHERE id=$1", id)
	if err != nil {
		c.JSON(404, gin.H{"message": "note not found"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"message": "note deleted successfully"})
}

func DeleteByDoctorID(c *gin.Context) {
	uuidParam := c.Param("uuid")

	_, err := databaseConn.Exec("DELETE FROM timetable WHERE doctor_uuid=$1", uuidParam)
	if err != nil {
		c.JSON(404, gin.H{"message": "note not found"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"message": "note deleted successfully"})
}

func DeleteByHospitalID(c *gin.Context) {
	uuidParam := c.Param("uuid")

	_, err := databaseConn.Exec("DELETE FROM timetable WHERE hospital_uuid=$1", uuidParam)
	if err != nil {
		c.JSON(404, gin.H{"message": "note not found"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"message": "note deleted successfully"})
}

func GetByHospitalUUID(c *gin.Context) {
	var allNotes []models.FullTimetable
	uuid := c.Param("uuid")
	token, err := c.Cookie("tokenAccess")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	fromParam, toParam := c.Query("from"), c.Query("to")

	if !_IsParamsValid(fromParam, toParam) {
		c.JSON(400, gin.H{"message": "not ISO8601 parameters"})
		return
	}

	if !_IsHospitalExist(uuid, token) {
		c.JSON(404, gin.H{"error": "hospital not found"})
	}

	rows, err := databaseConn.Query(`SELECT * FROM timetable WHERE hospital_uuid=$1`, uuid)
	if err != nil {
		c.JSON(501, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	for rows.Next() {
		note := models.FullTimetable{}
		err := rows.Scan(&note.ID, &note.HospitalUUID, &note.DoctorUUID, &note.From, &note.To, &note.Room)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		allNotes = append(allNotes, note)
	}

	c.JSON(200, allNotes)
}

func GetByDoctorUUID(c *gin.Context) {
	var allNotes []models.FullTimetable
	uuid := c.Param("uuid")
	token, err := c.Cookie("tokenAccess")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	fromParam, toParam := c.Query("from"), c.Query("to")

	if !_IsParamsValid(fromParam, toParam) {
		c.JSON(400, gin.H{"message": "not ISO8601 parameters"})
		return
	}

	if !_IsDoctorExist(uuid, token) {
		c.JSON(404, gin.H{"error": "doctor not found"})
		return
	}

	rows, err := databaseConn.Query(`SELECT * FROM timetable WHERE doctor_uuid=$1`, uuid)
	if err != nil {
		c.JSON(501, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	for rows.Next() {
		note := models.FullTimetable{}
		err := rows.Scan(&note.ID, &note.HospitalUUID, &note.DoctorUUID, &note.From, &note.To, &note.Room)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		allNotes = append(allNotes, note)
	}

	c.JSON(200, allNotes)
}

func GetByHospitalUUIDAndRoom(c *gin.Context) {
	var allNotes []models.FullTimetable
	uuid, room := c.Param("uuid"), c.Param("room")
	token, err := c.Cookie("tokenAccess")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	fromParam, toParam := c.Query("from"), c.Query("to")

	if !_IsParamsValid(fromParam, toParam) {
		c.JSON(400, gin.H{"message": "not ISO8601 parameters"})
		return
	}

	if !_IsHospitalExist(uuid, token) {
		c.JSON(404, gin.H{"error": "hospital not found"})
		return
	}

	if !_IsRoomExist(uuid, token, room) {
		c.JSON(404, gin.H{"error": "room not found"})
		return
	}

	rows, err := databaseConn.Query(`SELECT * FROM timetable WHERE hospital_uuid=$1 AND room=$2`, uuid, room)
	if err != nil {
		c.JSON(501, gin.H{"message": err.Error()})
		return
	}
	
	if !rows.Next() {
		c.JSON(404, gin.H{"message": "timetable not ready"})
		return
	}

	for rows.Next() {
		note := models.FullTimetable{}
		err := rows.Scan(&note.ID, &note.HospitalUUID, &note.DoctorUUID, &note.From, &note.To, &note.Room)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		allNotes = append(allNotes, note)
	}

	c.JSON(200, allNotes)
}