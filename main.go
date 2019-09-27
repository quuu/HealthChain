package main

import ( 
	"fmt"
	"flag"
	)




func main(){
	flag.Parse()

	fmt.Println("Starting discovery")

	go discovery()

  api()

	
	return 
}
