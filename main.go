package main

import (
	"flag"
	"fmt"
)

func main() {
	flag.Parse()

	fmt.Println("Starting discovery")

	//pd := CreatePeerDriver()
	//pd.Discovery()

	api := NewAPI()

	api.Run()

	return
}
