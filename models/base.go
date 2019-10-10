package models

import (
	u "HealthChain/utils"
	"encoding/json"
	f "fmt"

	"github.com/asdine/storm"
)

// dont think i need this
// var db *storm.DB

type Patient struct {
	PatientKey string   `storm:"id"` // public key to access a patients records
	Records    []string // encrypted json strings for any record associated with this pateint
	Node       string   // identifies what node this record is on
}

func InitDB(init_file string) {
	f.Println("Initializing db")
	db := GetDB()
	if init_file == "" {
		f.Println("Nothing to initalize")
		db.Close()
		return
	}
	defer db.Close()

	// todo
	// change logic below to test the base funcs
	rows := u.Mapsfjson(init_file)
	f.Println("Obtained records from json file")
	for _, row := range rows {
		patient_id := DB_key(row)
		var records []string
		jstring, _ := json.Marshal(row)
		encrypted_jstr := u.Encrypt_json_string(string(jstring), patient_id)
		records = append(records, string(encrypted_jstr))
		p := Patient{PatientKey: patient_id, Records: records, Node: "hc_1"}
		f.Println("Saving -> ", p)
		save_err := db.Set("patients", patient_id, &p)
		if save_err != nil {
			f.Println(save_err.Error())
		}
	}
	var test Patient
	get_err := db.Get("patients", "Laurie Feliciano6361781291938-11-18US", &test)
	if get_err != nil {
		f.Print(get_err.Error())
	}
	f.Println("This should be a record :) -> ", test)

}

func DB_key(patient_record u.Json) string {
	fname := u.Valtostr(patient_record["full_name"])
	ssn := u.Valtostr(patient_record["ssn"])
	dob := u.Valtostr(patient_record["dob"])
	country := u.Valtostr(patient_record["Country"])
	return fname + ssn + dob + country
}

func GetDB() *storm.DB {
	db, err := storm.Open("hc.db")
	if err != nil {
		panic(err)
	}
	return db
}

func (patient *Patient) Validate() (u.Json, bool) {
	// validate that Patient json is well formed
	return u.Message(true, "success"), true
}

func GetPatient(key string) *Patient {
	db := GetDB()

	patient := &Patient{}
	err := db.Get("patients", key, &patient)
	if err != nil {
		return nil // handler must catch if the get query returns nil
	}
	return patient
}

func (patient *Patient) AddPatient() u.Json {
	if resp, ok := patient.Validate(); !ok {
		return resp
	}

	db := GetDB()
	err := db.Set("patients", patient.PatientKey, &patient)
	if err != nil {
		f.Println(err)
		return u.Message(false, err.Error())
	}
	resp := u.Message(true, "success")
	resp["patient"] = patient
	return resp

}
