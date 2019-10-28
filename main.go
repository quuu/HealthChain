package main

import (
	"flag"
	"fmt"

	"github.com/asdine/storm/v3"
)

func main() {
	flag.Parse()

	fmt.Println("Starting discovery")

	//pd := CreatePeerDriver()
	//pd.Discovery()

	db, _ := storm.Open("my.db")
	api := NewAPI(db)

	api.Run()

	return
}
