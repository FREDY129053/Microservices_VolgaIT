package models

import (
	"time"
)

type HistoryInfo struct {
	Date         time.Time `json:"date"`
	PatientUUID  string    `json:"patient_uuid"`
	HospitalUUID string    `json:"hospital_uuid"`
	DoctorUUID   string    `json:"doctor_uuid"`
	Room         string    `json:"room"`
	Data         string    `json:"data"`
}

type HospitalShortInfo struct {
	UUID string `json:"hospital_uuid"`
	Name string `json:"name"`
}

type DoctorShortInfo struct {
	UUID      string `json:"doctor_uuid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type PatientShortInfo struct {
	UUID      string `json:"patient_uuid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type FullHistory struct {
	ID           int       `json:"id"`
	Date         time.Time `json:"date"`
	PatientUUID  string    `json:"patient_uuid"`
	HospitalUUID string    `json:"hospital_uuid"`
	DoctorUUID   string    `json:"doctor_uuid"`
	Room         string    `json:"room"`
	Data         string    `json:"data"`
}
