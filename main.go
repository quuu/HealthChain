package main

import (
	"flag"
	"fmt"
)

func main() {
	flag.Parse()

	fmt.Println("Starting discovery")

	// open the database

	// use the same database object for the peer driver
	pd := CreatePeerDriver()
	go pd.Discovery()

	// as well as the api
	api := NewAPI(pd.uuid)
	api.Run()

	return
}
