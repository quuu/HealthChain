package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
)

// API
type API struct {
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

func (a *API) NewRecordHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	for _, element := range r.Form {
		log.Println(element)
	}
}

func NewAPI() *API {

	// initialize a new api
	a := &API{}

	log.Printf("initializing api \n")

	r := chi.NewRouter()

	// create the database to reference
	db, _ := storm.Open("my.db")

	a.store = db

	rec := Record{ID: "someoneelse", Message: []byte("testing"), Date: time.Now()}

	err := db.Save(&rec)
	if err != nil {
		log.Printf("errored at %s", err)
	}

	rec2 := Record{ID: "me", Message: []byte("asdf"), Date: time.Now()}
	err = db.Save(&rec2)
	if err != nil {
		log.Printf("errored at %s", err)
	}

	r.Use(middleware.DefaultCompress)

	// r.Method("GET", "/static/*", http.FileServer)

	r.Get("/", a.IndexHandler)

	r.Get("/all_records", a.AllRecordsHandler)

	r.Get("/peers", a.PeersHandler)

	r.Post("/new_record", a.NewRecordHandler)

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
