package main

import (
	"encoding/json"
	"net/http"

	"github.com/asdine/storm/v3"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
)

// API
type API struct {
	str string
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent("a", "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(b)
}

func PeersHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent("a", "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		// return err
	}
	w.Write(b)
}
func AllRecordsHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent("a", "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		// return err
	}
	w.Write(b)
}

func NewRecordHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	for _, element := range r.Form {
		log.Println(element)
	}
}
func api() {

	r := chi.NewRouter()

	db, _ := storm.Open("my.db")

	defer db.Close()

	r.Use(middleware.DefaultCompress)

	// r.Method("GET", "/static/*", http.FileServer)

	r.Get("/", IndexHandler)

	r.Get("/all_reocrds", AllRecordsHandler)

	r.Get("/peers", PeersHandler)

	r.Post("/new_record", NewRecordHandler)

	err := http.ListenAndServe(":3000", r)
	if err != nil {

		log.WithError(err).Error("unable to listen and serve")
	}

}
