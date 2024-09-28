package models

type HospitalInfo struct {
	UUID         string   `json:"uuid"`
	Name         string   `json:"name"`
	Address      string   `json:"address"`
	ContactPhone string   `json:"contact_phone"`
	Rooms        []string `json:"rooms,omitempty"`
}

type AddHospitalInfo struct {
	Name         string   `json:"name"`
	Address      string   `json:"address"`
	ContactPhone string   `json:"contact_phone"`
	Rooms        []string `json:"rooms"`
}