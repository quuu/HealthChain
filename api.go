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
	uuid "github.com/satori/go.uuid"
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

func (a *API) DecryptPatientRecords(first string, last string, country string, code string) {

	enc := a.FetchEncrypted()

	log.Println("fetching all records")
	hash_key := GetHash(first, last, country, code)
	for _, elem := range enc {

		log.Println("this is a record")
		log.Println(elem.Contents)
		log.Println("attempting to decrypt")
		temp := Decrypt(hash_key, elem.Contents)
		log.Println(temp)
	}

}
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
	r.ParseForm()

	first := r.Form.Get("first")
	last := r.Form.Get("last")
	country := r.Form.Get("country")
	code := r.Form.Get("code")

	hash_key := GetHash(first, last, country, code)
	log.Println("got hash key  ")
	log.Println(hash_key)

	enc := a.FetchEncrypted()
	for _, record := range enc {

		log.Println("this is a record")
		log.Println(record.Contents)
		log.Println("attempting to decrypt")
		temp := Decrypt(hash_key, record.Contents)
		log.Println(string(temp))
	}

	// var enc []*EncryptedRecord

	// err := a.store.Find("ID", hash_key, &enc)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// fmt.Println("got decoded")
	// fmt.Println(enc)

	// w.Write([]byte("found stuff"))

	// a.DecryptPatientRecords(first, last, country, code)
}

func (a *API) StoreRecord(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	first := r.Form.Get("first")
	last := r.Form.Get("last")
	country := r.Form.Get("country")
	code := r.Form.Get("code")

	// get the hash of the user
	hash_key := GetHash(first, last, country, code)

	log.Println("got the hash key " + string(hash_key))
	//get the messaage
	appointment_info := r.Form.Get("appointment_info")

	// create a new record
	rec := &Record{
		Message: appointment_info,
		Date:    time.Now(),
	}

	log.Println("the mesage is ")
	log.Println(rec.Message)
	// log.Println(string(rec.Message))

	b, err := json.Marshal(rec)
	if err != nil {
		panic(err.Error())
	}

	encrypted := Encrypt(hash_key, b)

	// save the encrypted record with the id
	// being the hash key
	// TODO
	// find a better way to store the id for lookup
	enc := EncryptedRecord{
		ID:       uuid.NewV4().String(),
		Contents: encrypted,
	}

	// actually save it into the database
	err = a.store.Save(&enc)
	if err != nil {
		panic(err.Error())
	}

	//temporary code to make sure the encryption
	// can be decrypted
	// decrypted := Decrypt(hash_key, encrypted)
	// log.Println(string(decrypted))

	// var re *Record
	// json.Unmarshal(decrypted, &re)
	// log.Println("decrypted message is " + string(re.Message))

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
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

// function that will return the md5 hash encoding given 4 string parameters
func GetHash(first string, last string, country string, code string) []byte {

	hasher := sha256.New()
	hasher.Write([]byte(first + last + country + code))

	log.Println("made the hash " + string(hasher.Sum(nil)))
	return hasher.Sum(nil)
	// gcm, err := cipher.NewGCM(block)
	// if err != nil {
	// 	panic(err.Error())
	// }

	// nonce := make([]byte, gcm.NonceSize())
	// if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
	// 	panic(err.Error())
	// }
	// cipherText := gcm.Seal(nonce, nonce, )
}
func (a *API) NewRecordHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	for _, element := range r.Form {
		log.Println(element)
	}
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

	rec := Record{ID: "someoneelse", Message: "testing", Date: time.Now()}

	err := store.Save(&rec)
	if err != nil {
		log.Printf("errored at %s", err)
	}

	rec2 := Record{ID: "me", Message: "asdf", Date: time.Now()}
	err = store.Save(&rec2)
	if err != nil {
		log.Printf("errored at %s", err)
	}

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
