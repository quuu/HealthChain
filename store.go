package main

import (
	"time"
)

type FormData struct {
	First            string `json:"first"`
	Last             string `json:"last"`
	Country          string `json:"country"`
	Code             string `json:"code"`
	Appointment_Info struct {
		Summary     string `json:"summary"`
		Height      string `json:"height"`
		Weight      string `json:"weight"`
		Vaccination string `json:"vaccination"`
		Sickness    string `json:"sickness"`
		Eyesight    string `json:"eyesight"`
		Date        time.Time
	} `json:"appointment_info"`
}

type Record struct {
	ID      string    `json:"ID"`
	Message []byte    `json:"Message"`
	Date    time.Time `json:"Date"`
	Type    string    `json:"Type"`
}

type Patient struct {
	PatientKey string   `storm:"id"` // public key to access a paitnets records
	Records    []Record // encrypted json string for records asspcotaetd with this patient
	Node       string   // identifies what node this record patient was created on
}

type EncryptedRecord struct {
	PatientID string
	Contents  []byte `storm:"id"`
}

type DecryptedRecord struct {
	ID       string
	Contents []byte
}
