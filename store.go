package main

import (
	"time"
)

type Record struct {
	ID      string
	Message []byte
	Date    time.Time
}
