package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	http.Handle("/", http.FileServer(http.Dir("./client")))

	log.Fatal(http.ListenAndServe(":4000", nil))

	fmt.Println("Server up on port 4000")

}
