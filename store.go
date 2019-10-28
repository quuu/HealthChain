package main

import (
	"time"
)

type Record struct {
	ID      string
	Message []byte
	Date    time.Time
}

type EncryptedRecord struct {
	ID       string
	Contents []byte
}
