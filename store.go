package main

import (
	"time"
)
// FormData
// json template for patient record
// When recading a record.json only reads from fields in struct 
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

// Record
// Encapsulated Record with unique id, date.time 
type Record struct {
	ID      []byte    `json:"ID"`
	Message []byte    `json:"Message"`
	Date    time.Time `json:"Date"`
	Type    string    `json:"Type"`
}

// Patient
// Encapsulated all records associated with PatientKey such that all patientKeys 
// are unique and is hashed from the patients credentials.
// patientkey is obtained from encryption hash.
type Patient struct {
	PatientKey []byte   `storm:"id"` // public key to access a paitnets records
	Records    []Record // encrypted json string for records asspcotaetd with this patient
	Node       string   // identifies what node this record patient was created on
}

// EncryptedRecord
type EncryptedRecord struct {
	Contents  []byte `storm:"id"`
	PatientID []byte
}

//DecryptedRecord
type DecryptedRecord struct {
	ID       string
	Contents []byte
}
