package models

import (
	"time"
)

type Timetable struct {
	HospitalUUID string    `json:"hospital_uuid"`
	DoctorUUID   string    `json:"doctor_uuid"`
	From         time.Time `json:"from"`
	To           time.Time `json:"to"`
	Room         string    `json:"room"`
}

type FullTimetable struct {
	ID int `json:"uuid"`
	Timetable
}