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
type API struct {
	me    string
	str   string
	store *storm.DB
	r     http.Handler
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

	var encrypted []*EncryptedRecord
	err := a.store.All(&encrypted)

	if err != nil {
		panic(err.Error())
	}

	return encrypted
}

// handler to get all encrypted records currently being stored
func (a *API) AllRecordsHandler(w http.ResponseWriter, r *http.Request) {

	var records []*EncryptedRecord
	err := a.store.All(&records)

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

	// TODO

	// get hash_key index in records table

	// unhash with the key

	// unhash each appointment

	// return json

	// OVERLY SIMPLIFIED

	var records []string
	enc := a.FetchEncrypted()
	for _, record := range enc {

		// if the decryption was successful
		temp := Decrypt(hash_key, record.Contents)
		if temp != nil {
			records = append(records, string(temp))
		}
	}

	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent(records, "", "")
	if err != nil {
		panic(err.Error())
	}
	w.Write(b)

}

func (a *API) StoreRecord(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(0)

	first := r.FormValue("first")
	last := r.FormValue("last")
	country := r.FormValue("country")
	code := r.FormValue("code")

	// get the hash of the user
	hash_key := GetHash(first, last, country, code)

	//get the messaage
	appointment_info := r.Form.Get("appointment_info")

	log.Println("appointment info " + appointment_info)

	// TODO

	//hash appointment_info with hash_key

	//fetch all appointment's of the user currently

	//unhash the contents of the user

	// access the list that stores the hashed appointments

	// append the hashed appointment_info to the appointments

	// no longer using randomly generated UUID
	// unique_id := uuid.NewV4().String()
	// create a new record
	rec := &Record{
		ID:      string(hash_key),
		Message: appointment_info,
		Date:    time.Now(),
	}

	// conver the contents to a byte array
	b, err := json.Marshal(rec)
	if err != nil {
		panic(err.Error())
	}

	// encrypt the record
	encrypted := Encrypt(hash_key, b)

	// save the encrypted record with the id
	// being the hash key
	// TODO
	// find a better way to store the id for lookup

	// using a unique_id to prevent duplicate saves
	enc := EncryptedRecord{
		ID:       string(hash_key),
		Contents: encrypted,
	}

	// actually save it into the database
	err = a.store.Save(&enc)
	if err != nil {
		panic(err.Error())
	}
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

func NewAPI(store *storm.DB, uuid string) *API {

	// initialize a new api
	a := &API{
		me:    uuid,
		store: store,
	}

	log.Printf("initializing api \n")

	r := chi.NewRouter()

	// create the database to reference
	// rec := Record{ID: "someoneelse", Message: "testing", Date: time.Now()}

	// err := store.Save(&rec)
	// if err != nil {
	// 	log.Printf("errored at %s", err)
	// }

	// rec2 := Record{ID: "me", Message: "asdf", Date: time.Now()}
	// err = store.Save(&rec2)
	// if err != nil {
	// 	log.Printf("errored at %s", err)
	// }

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
	defer a.store.Close()

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

	patient := &Patient{}
	err := db.Get("patients", key, &patient)
	if err != nil {
		return nil // handler must catch if the get query returns nil
	}
	return patient
}

func (patient *Patient) AddPatient() map[string]interface{} {

	db := GetDB()
	err := db.Set("patients", patient.PatientKey, &patient)
	if err != nil {
		log.Println(err)
		return Message(false, err.Error())
	}
	resp := Message(true, "success")
	resp["patient"] = patient
	return resp

}

func AddRecord(key string, record Record) map[string]interface{} {
	db := GetDB()
	patient := &Patient{}
	err := db.Get("patients", key, &patient)
	if err != nil {
		log.Println("The Patient could not be found")
		return Message(false, err.Error())
	}

	var recs []Record

	recs = patient.Records
	recs = append(recs, record)
	patient.Records = recs
	err_set := db.Set("patients", patient.PatientKey, &patient)
	if err_set != nil {
		log.Println(err)
		return Message(false, err.Error())
	}
	resp := Message(true, "success")
	resp["patient"] = patient
	return resp

}
