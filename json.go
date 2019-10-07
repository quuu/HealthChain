package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type M map[string]interface{}

func mapsfjson() []M {
	//open the file
	jsonFile, err := os.Open("hc_data.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	// read bytes in file
	byteVal, _ := ioutil.ReadAll(jsonFile)

	// store the rows of json file into list of maps
	var rows []M
	json.Unmarshal([]byte(byteVal), &rows)

	return rows
}

func main() {
	rows := mapsfjson()
	fmt.Println(rows[0]["ssn"])
}
