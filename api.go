package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
)

// API
type API struct {
	me  string
	str string
	r   http.Handler
}

// handler for main page
// route it to index.html
func (a *API) IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent("a", "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(b)
}

// handler to find out which peers are currently connected
func (a *API) PeersHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent("a", "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		// return err
	}
	w.Write(b)
}

// gets all encrypted records
func (a *API) FetchEncrypted() []*EncryptedRecord {
	db := GetDB()
	defer db.Close()
	var encrypted []*EncryptedRecord
	err := db.All(&encrypted)

	if err != nil {
		panic(err.Error())
	}

	return encrypted
}

// handler to get all encrypted records currently being stored
func (a *API) AllRecordsHandler(w http.ResponseWriter, r *http.Request) {
	db := GetDB()
	defer db.Close()

	var records []*EncryptedRecord
	err := db.All(&records)

	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent(records, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		// return err
	}
	w.Write(b)
}

func (a *API) GetRecords(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(0)

	first := r.FormValue("first")
	last := r.FormValue("last")
	country := r.FormValue("country")
	code := r.FormValue("code")

	hash_key := GetHash(first, last, country, code)

	patient := &Patient{}
	patient = GetPatient(string(hash_key))
	if patient == nil {
		log.Println("No records for this patient")
		return
	}
	log.Println("Found patient", patient)

	// TODO

	// get hash_key index in records table

	// unhash with the key

	// unhash each appointment

	// return json

	// OVERLY SIMPLIFIED

	// todo
	// case where the key does not hash to a patient

	var records []Record
	records = patient.Records

	var decrypted_records []string

	// var records_decrypt []string

	for _, record := range records {

		fmt.Println("drcypring ")
		decrypted := Decrypt(hash_key, record.Message)
		fmt.Println(record.Date)
		decrypted_records = append(decrypted_records, string(decrypted))

		// err := json.Unmarshal(record.Message, &record.Message)
		// if err != nil {
		// 	log.Println("something wrong with unmasrhslling this json")
		// }
		// if the decryption was successful
		// temp := Decrypt(hash_key, record.Message)
		// if temp != nil {
		// 	err := json.Unmarshal(record, &record.Message)
		// 	if err != nil {
		// 		log.Println("something wrong with unmasrhslling this json")
		// 	}
		// 	// records_decrypt = append(records_decrypt, string(temp))
		// }
	}

	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent(decrypted_records, "", "")
	if err != nil {
		panic(err.Error())
	}
	w.Write(b)

}

func (a *API) StoreRecord(w http.ResponseWriter, r *http.Request) {

	// decode the response body
	decoder := json.NewDecoder(r.Body)

	// extract the data into FormData struct
	var response FormData

	err := decoder.Decode(&response)
	if err != nil {
		panic(err)
	}

	response.Appointment_Info.Date = time.Now()

	// get the appointment info
	output, err := json.Marshal(response.Appointment_Info)
	if err != nil {

	}
	fmt.Println(string(output))

	first := response.First
	last := response.Last
	country := response.Country
	code := response.Code

	// get the hash of the user
	hash_key := GetHash(first, last, country, code)

	//get the messaage
	apt_encrypt := Encrypt(hash_key, output)

	// apt_json_encyp := Encrypt(hash_key, apt_json)

	//rec_tostore := Record{ID: string(hash_key), Message: apt_json_encyp, Date: time.Now()}
	rec_to_store := Record{ID: string(hash_key), Message: apt_encrypt, Date: time.Now()}

	// no longer using randomly generated UUID
	// unique_id := uuid.NewV4().String()

	// REVIEW: whats wrong with the record having a rand gen uuid?
	//		the records are encapsulated within patient data

	// gets patient if already present else makes a new patient to store
	p := &Patient{}
	var records_temp []Record
	p = GetPatient(string(hash_key)) // NOTE this is a patch, method should take a []byte
	if p == nil {
		log.Println("Patient does not exsist, making new patient")
		p := Patient{PatientKey: string(hash_key), Records: records_temp, Node: "hc_1"}
		AddPatient(p)
	}

	// actually save it into the database
	p.AddRecord(string(hash_key), rec_to_store)
	// if err != nil {
	// 	panic(err.Error())
	// }
	w.Write([]byte("Saved!"))
}

// function that will return an encrypted record given
// a hashed encoding and a stream of bytes to encode
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
	cipher_text := gcm.Seal(nonce, nonce, data, nil)
	return cipher_text
}

// function that will return a decrypted record given
// a hashed encoding and a stream of bytes that has been encoded
// TODO
// if the result is not fully decoded, return something meaningful
func Decrypt(encoding []byte, data []byte) []byte {
	block, _ := aes.NewCipher(encoding)
	gcm, err := cipher.NewGCM((block))
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, cipher_text := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, cipher_text, nil)

	// failed to decrypt, would normally throw an error
	if err != nil {
		// panic(err.Error())
		return nil
	}
	return plaintext
}

// function that will return the md5 hash encoding given 4 string parameters
func GetHash(first string, last string, country string, code string) []byte {

	// utilize sha256 to fix the hash to be 32 bytes long
	hasher := sha256.New()
	hasher.Write([]byte(first + last + country + code))

	log.Println("made the hash " + string(hasher.Sum(nil)))
	return hasher.Sum(nil)
}

func NewAPI(uuid string) *API {

	// initialize a new api
	a := &API{
		me: uuid,
	}

	log.Printf("initializing api \n")

	r := chi.NewRouter()

	r.Use(middleware.DefaultCompress)

	// r.Method("GET", "/static/*", http.FileServer)

	r.Get("/", a.IndexHandler)

	r.Get("/all_records", a.AllRecordsHandler)

	r.Get("/peers", a.PeersHandler)

	r.Post("/new_record", a.StoreRecord)

	r.Post("/get_records", a.GetRecords)

	a.r = r

	return a
}

func (a *API) Run() {
	// close the store once listening is done
	// listen for requests
	log.Printf("listening")
	err := http.ListenAndServe(":3000", a.r)
	if err != nil {

		log.WithError(err).Error("unable to listen and serve")
	}

}

// ================== RICH/methods ==================
// I needed to do a little bit of refactoring

// message for sending to http client
func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func GetDB() *storm.DB {
	db, err := storm.Open("hc.db")
	if err != nil {
		panic(err)
	}
	return db
}

func GetPatient(key string) *Patient {
	db := GetDB()
	defer db.Close()

	patient := &Patient{}
	err := db.Get("patients", key, &patient)
	if err != nil {
		return nil // han dler must catch if the get query returns nil
	}
	return patient
}

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

func (patient *Patient) AddRecord(key string, record Record) map[string]interface{} {
	db := GetDB()
	defer db.Close()

	err := db.Get("patients", key, &patient)
	if err != nil {
		log.Println("The Patient could not be found")
		return Message(false, err.Error())
	}

	// var recs []Record

	// recs = patient.Records
	// recs = append(recs, record)
	patient.Records = append(patient.Records, record)
	err_set := db.Set("patients", key, &patient)
	if err_set != nil {
		log.Println(err)
		return Message(false, err.Error())
	}
	resp := Message(true, "success")
	resp["patient"] = patient
	return resp

}
