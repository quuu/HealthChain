package main

import (
	"flag"
	"fmt"

	"github.com/asdine/storm/v3"
)

func main() {
	flag.Parse()

	fmt.Println("Starting discovery")

	// open the database
	db, _ := storm.Open("my.db")

	// use the same database object for the peer driver
	pd := CreatePeerDriver(db)
	go pd.Discovery()

	// as well as the api
	api := NewAPI(pd.uuid, db)
	api.Run()

	return
}
