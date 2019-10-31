package main

import (
	"time"
)

type Record struct {
	ID      string    `json:"ID"`
	Message string    `json:"Message"`
	Date    time.Time `json:"Date"`
}

type EncryptedRecord struct {
	ID       string
	Contents []byte
}

type DecryptedRecord struct {
	ID       string
	Contents []byte
}
