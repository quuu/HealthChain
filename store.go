package main

import (
	"time"
)

type Record struct {
	ID      string    `json:"ID"`
	Message []byte    `json:"Message"`
	Date    time.Time `json:"Date"`
}

type Patient struct {
	PatientKey string   `storm:"id"` // public key to access a paitnets records
	Records    []Record // encrypted json string for records asspcotaetd with this patient
	Node       string   // identifies what node this record patient was created on
}

type EncryptedRecord struct {
	ID       string
	Contents []byte
}

type DecryptedRecord struct {
	ID       string
	Contents []byte
}
