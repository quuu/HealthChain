package main

import (
	"time"
)

type Record struct {
	ID      string
	Message string
	Date    time.Time
}

type EncryptedRecord struct {
	ID       string
	Contents []byte
}

type DecryptedRecord struct {
	ID       string
	Contents []byte
}
