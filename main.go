package main

import (
	"flag"
	"fmt"
)

func main() {
	flag.Parse()

	fmt.Println("Starting discovery")

	discovery()

	api()

	return
}
