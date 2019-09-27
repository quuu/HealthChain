package main

import (
	"flag"
	"fmt"
)

func main() {
	flag.Parse()

	fmt.Println("Starting discovery")

	go discovery()

	api()

	return
}
