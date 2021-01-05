package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
)

// API
// me: unique identifier of the node
// str: self name, to be used later
// r: http.Request handler
type API struct {
	me  string
	str string
	r   http.Handler
}

// IndexHandler
// handler for main page routed at index.html
func (a *API) IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent("Index", "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(b)
}

// PeersHandler
// Handler to find out which peers are currently connected
// PeersHandler is called to check if there are updates to be fetched within the
//  network
func (a *API) PeersHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent("a", "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(b)
}

// AllRecordsHandler
// Writes all encrypted records currently located in the local db to responceWriter
// Does not modify the contents of the db, but is called when syncing peers
func (a *API) AllRecordsHandler(w http.ResponseWriter, r *http.Request) {
	db := GetDB()
	defer db.Close()

	var records []*EncryptedRecord
	err := db.All(&records)

	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent(records, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(b)
}

// GetRecords
// Writes/Returns records from patient p
// Requires p's first/last name, country_code and unique_code (ssn)
//  and patient existing in the local db.
// Look up in DB is O(1)
func (a *API) GetRecords(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(0)

	first := r.FormValue("first")
	last := r.FormValue("last")
	country := r.FormValue("country")
	code := r.FormValue("code")

	hashKey := GetHash(first, last, country, code)

	patient := &Patient{}
	patient = GetPatient(hashKey)

	// Case where the coressponding patient does not locally exsist.
	// Returns nil if records are not present.
	if patient == nil {
		log.Printf("No records for this patient % x\n", hashKey)
		return
	} else {
		log.Printf("Found patient% x\n", hashKey)
	}

	var records []Record
	records = patient.Records

	var decryptedRecprds []string

	for _, record := range records {

		// decrypt the record
		decrypted := Decrypt(hashKey, record.Message)

		// add it to the return list
		decryptedRecprds = append(decryptedRecprds, string(decrypted))
	}

	w.Header().Set("Content-Type", "application/json")
	decrytpedRecordsJson, err := json.MarshalIndent(decryptedRecprds, "", "")
	if err != nil {
		panic(err.Error())
	}
	w.Write(decrytpedRecordsJson) // returns to ResponseWriter

}

// StoreRecord
// Saves encyrpted record locally to db
// Requires patients creditials
func (a *API) StoreRecord(w http.ResponseWriter, r *http.Request) {

	// decode the response body
	decoder := json.NewDecoder(r.Body)

	// extract the data into FormData struct
	var response FormData

	err := decoder.Decode(&response)
	if err != nil {
		panic(err)
	}

	response.AppointmentInfo.Date = time.Now()

	// get the appointment info
	output, err := json.Marshal(response.AppointmentInfo)
	if err != nil {

	}

	first := response.First
	last := response.Last
	country := response.Country
	code := response.Code

	// get the hash of the patient
	hashKey := GetHash(first, last, country, code)

	// get the messaage
	encryptedAppointment := Encrypt(hashKey, output)

	recordToStore := Record{ID: hashKey, Message: encryptedAppointment, Date: time.Now(), Type: "Message"}

	// gets patient if already present else makes a new patient struct to store
	p := &Patient{}
	var tempRecords []Record
	p = GetPatient(hashKey)

	if p == nil {
		log.Println("Patient does not exsist, making new patient")
		p := Patient{PatientKey: hashKey, Records: tempRecords, Node: "hc_1"}
		AddPatient(p)
	}

	enc := &EncryptedRecord{PatientID: hashKey, Contents: recordToStore.Message}
	peerDB := PublicDB()

	// peer database usage
	err = peerDB.Save(enc)
	if err != nil {
		panic(err)
	}

	peerDB.Close()

	// actually save it into the database
	p.AddRecord(hashKey, recordToStore)
	w.Write([]byte("Saved!")) // success return message to ResponseWriter
}

// Encrypt
// Return an encrypted record (data) and encrypts it with encoding
// Data can only be decrypted via the same encoding
func Encrypt(encoding []byte, data []byte) []byte {
	block, _ := aes.NewCipher(encoding)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	cipherText := gcm.Seal(nonce, nonce, data, nil)
	return cipherText
}

// Decrypt
// Returns decrypted data from encoding
// Decrypted data encrypted by Encrypt()
func Decrypt(encoding []byte, data []byte) []byte {
	block, _ := aes.NewCipher(encoding)
	gcm, err := cipher.NewGCM((block))
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, cipherText := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, cipherText, nil)

	// failed to decrypt, would normally throw an error
	if err != nil {
		return nil
	}
	return plaintext
}

// GetHash
// Returns the md5 hash encoding given 4 string parameters that represent patient
//  credientials
func GetHash(first string, last string, country string, code string) []byte {

	// utilize sha256 to fix the hash to be 32 bytes long
	hasher := sha256.New()
	hasher.Write([]byte(first + last + country + code))
	return hasher.Sum(nil)
}

// NewAPI
// Returns a API stucct that wraps chi router with function handlers
func NewAPI(uuid string) *API {

	// initialize a new api
	a := &API{
		me: uuid,
	}

	log.Printf("initializing api \n")

	r := chi.NewRouter() // router

	r.Use(middleware.DefaultCompress) // middleware

	// handlers
	r.Get("/", a.IndexHandler)

	r.Get("/all_records", a.AllRecordsHandler)

	r.Get("/peers", a.PeersHandler)

	r.Post("/new_record", a.StoreRecord)

	r.Post("/get_records", a.GetRecords)

	a.r = r

	return a
}

// Run
// Modifies API router to listen and Server on port :3000
func (a *API) Run() {
	// close the store once listening is done
	// listen for requests
	log.Printf("listening")
	err := http.ListenAndServe(":3000", a.r)
	if err != nil {

		log.WithError(err).Error("unable to listen and serve")
	}

}

// Message
// json for sending information/logging to  http client
// Return map[string]interface{} with status and message fields
func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}


// GetDB
// Returns refernece to local db for addition
// Used internally for fast lookup of records
/*
	DESGIN PATTERN: SINGLETON
	- Singleton is used for two reasons
		1. the file location can be abstracted away for clients need to 
			acess the db s.t. the client does not need to change code
			 for calling access to the db
		2. Helps regulate mutex for db - see storm documentation
*/
func GetDB() *storm.DB {
	db, err := storm.Open("hc.db")
	if err != nil {
		panic(err)
	}
	return db
}

// PublicDB
// Returns reference to public db for addition
// Used when records are synced across peers
func PublicDB() *storm.DB {
	db, err := storm.Open("records.db")
	if err != nil {
		panic(err)
	}
	return db
}

// GetPatient
// Returns patient struct of patient that contains encrypted
//  records associtated with key
func GetPatient(key []byte) *Patient {
	db := GetDB()
	defer db.Close()

	patient := &Patient{}
	err := db.Get("patients", key, &patient)
	if err != nil {
		return nil // han dler must catch if the get query returns nil
	}
	return patient
}

// AddPatient
// Modifies local db, adds patient to db
// Returns Message intended to be send to http.client
func AddPatient(patient Patient) map[string]interface{} {

	db := GetDB()
	defer db.Close()

	err := db.Set("patients", patient.PatientKey, &patient)
	if err != nil {
		log.Println(err)
		return Message(false, err.Error())
	}
	resp := Message(true, "success")
	resp["patient"] = patient
	return resp

}

// AddRecord
// Modicfies patient and local db, adds record assosciated with patient that is
//	returned from key.
// Requires that patient exsist in local db
// Returns Message intended to be send to http.client
func (patient *Patient) AddRecord(key []byte, record Record) map[string]interface{} {
	db := GetDB()
	defer db.Close()

	err := db.Get("patients", key, &patient)
	if err != nil {
		log.Println("The Patient could not be found")
		return Message(false, err.Error())
	}
	patient.Records = append(patient.Records, record)
	err = db.Set("patients", key, &patient)
	if err != nil {
		log.Println(err)
		return Message(false, err.Error())
	}
	resp := Message(true, "success")
	resp["patient"] = patient
	return resp

}