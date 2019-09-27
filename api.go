package main

import (
	"encoding/json"
	"net/http"

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
		// return err
	}
	w.Write(b)
	// return nil
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

func api() {

	r := chi.NewRou

	r.Use(middleware.DefaultCompress)

	// r.Method("GET", "/static/*", http.FileServer)

	r.Get("/", IndexHandler)

	r.Get("/all_reocrds", AllRecordsHandler)

	r.Get("/peers", PeersHandler)

	err := http.ListenAndServe(":3000", r)
	if err != nil {

		log.WithError(err).Error("unable to listen and serve")
	}

}
