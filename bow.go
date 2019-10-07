package main

import (
	"fmt"
	"time"

	"github.com/zippoxer/bow"
)

type Record struct {
	Id      string `bow:"key"`
	Tags    []string
	Created time.Time
}

func main() {

	db, err := bow.Open("test")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	r1 := Record{
		Id:      "id1",
		Tags:    []string{"hello", "world"},
		Created: time.Now(),
	}

	err1 := db.Bucket("records").Put(r1)
	if err != nil {
		fmt.Println(err1)
	}

	var r2 Record
	err2 := db.Bucket("records").Get("id1", &r2)
	if err != nil {
		fmt.Println(err2)
	}
	fmt.Println(r2)
	fmt.Println(r2.Id)

}
