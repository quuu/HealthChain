package main

import (
	"crypto/sha1"
	"encoding/hex"
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

func valtostr(v interface{}) string {

	str := fmt.Sprintf("%v", v)
	return str
}

func main() {
	rows := mapsfjson()
	// fmt.Println(rows[0])
	jstring, err := json.Marshal(rows[0])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	jsonStr := string(jstring)
	fmt.Println("heres the json string ->\n", jsonStr, "\n")

	var record M
	json.Unmarshal([]byte(jsonStr), &record)
	fmt.Println("Here is the json string converted back to map->\n", record, "\n")
	// here is how you marsh/unmarsh json into string and out of strings

	// todo
	// encrypt the json string, then decrypt, turn back into map
	// var fname string

	fname := valtostr(record["full_name"])

	ssn := valtostr(record["ssn"])
	dob := valtostr(record["dob"])
	country := valtostr(record["Country"])

	patient_id := fname + ssn + dob + country
	fmt.Println("Here is the patient_id string -> ", patient_id)

	hash := sha1.New()

	hash.Write([]byte(patient_id))
	sha1_hash := hex.EncodeToString(hash.Sum(nil))

	fmt.Println("Hash of ", patient_id, " -> ", sha1_hash)

}
