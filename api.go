package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
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

func (a *API) IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent("a", "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(b)
}

func (a *API) PeersHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent("a", "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		// return err
	}
	w.Write(b)
}
func (a *API) AllRecordsHandler(w http.ResponseWriter, r *http.Request) {

	var records []*Record
	err := a.store.All(&records)

	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent(records, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		// return err
	}
	w.Write(b)
}

func (a *API) StoreRecord(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	first := r.Form.Get("first")
	last := r.Form.Get("last")
	country := r.Form.Get("country")
	code := r.Form.Get("code")

	// get the hash of the user
	hash_key := GetHash(first, last, country, code)

	//get the messaage
	appointment_info := r.Form.Get("appointment_info")

	// get the current date

	// log.Println(appointment_info)
	rec := &Record{
		Message: []byte(appointment_info),
		Date:    time.Now(),
	}

	// log.Println("this messagE " + string(rec.Message))

	b, err := json.Marshal(rec)
	if err != nil {
		panic(err.Error())
	}

	log.Println(string(b))
	encrypted := Encrypt(hash_key, b)

	enc := EncryptedRecord{
		ID:       string(encrypted),
		Contents: encrypted,
	}
	err = a.store.Save(&enc)
	if err != nil {
		panic(err.Error())
	}

	// log.Println(string(encrypted))

	decrypted := Decrypt(hash_key, encrypted)
	log.Println(string(decrypted))
	var re *Record
	json.Unmarshal(decrypted, &re)
	log.Println("decrypted message is " + string(re.Message))

}

func Encrypt(block cipher.Block, data []byte) []byte {
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

func Decrypt(block cipher.Block, data []byte) []byte {
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

func GetHash(first string, last string, country string, code string) cipher.Block {

	hasher := md5.New()
	hasher.Write([]byte(first + last + country + code))
	encoding := hex.EncodeToString(hasher.Sum(nil))

	block, _ := aes.NewCipher([]byte(encoding))

	return block
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

	rec := Record{ID: "someoneelse", Message: []byte("testing"), Date: time.Now()}

	err := store.Save(&rec)
	if err != nil {
		log.Printf("errored at %s", err)
	}

	rec2 := Record{ID: "me", Message: []byte("asdf"), Date: time.Now()}
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
