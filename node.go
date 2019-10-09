package main

import (
	// utils

	"fmt"
	"net/http"
	"strings"

	// "contexts"
	// "flag"
	// "net/http"
	// "math/rand"

	// "bytes"
	"encoding/json"

	//chi-chi
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	// "github.com/go-chi/docgen"
	// "github.com/go-chi/render"

	// storm
	"github.com/asdine/storm"
)

// notes on progress
// establish nosql way to store data

// rest
// post -> add document to db assigned with hash value
// get -> get all documents associated with hash value

var router *chi.Mux

// var db storm whatever

const (
	dbName = "go_strom_crud"
	dbPass = "12345"
	dbHost = "localhost"
	dbPort = "33033"
)

func start_server() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to HealthChain :)\n"))
	})

	http.ListenAndServe(":3333", r)
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

type Patient struct {
	PatientKey string   `storm:"id"` // public key to access a patients records
	Records    []string // encrypted json strings for any record associated with this pateint
	Node       string   // identifies what node this record is on
}

func init_db() {
	db, err := storm.Open("hc.db")
	if err != nil {
		fmt.Println("db.Open exception")
	}
	defer db.Close()

	rows := mapsfjson()
	for _, row := range rows {
		patient_id := db_key(row)
		var records []string
		jstring, _ := json.Marshal(row)
		records = append(records, string(jstring))
		p := Patient{PatientKey: patient_id, Records: records, Node: "hc_1"}
		db.Save(&p)
	}
	var test Patient
	get_err := db.One("PatientKey", "Laurie Feliciano6361781291938-11-18US", &test)
	if get_err != nil {
		fmt.Print(get_err.Error())
	}
	fmt.Println("This should be a record :) -> ", test)

}
func main() {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to HealthChain :)\n"))
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	})
	init_db()
	// http.ListenAndServe(":3333", r)

}