package main

import (
	// utils
	"encoding/json"
	f "fmt"
	"net/http"
	"os"

	//chi-chi
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	// storm
	"github.com/asdine/storm"
	//hc
	"HealthChain/models"
)

// todo
// comments
// get comments working
// layout project the right way

var router *chi.Mux

func start_server() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/api/patient/new", func(w http.ResponseWriter, r *http.Request) {
		// call the func here ?
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to HealthChain :)\n"))
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	})

	http.ListenAndServe(":3333", r)
}

func view_db(db *storm.DB) {

	stats, err := json.Marshal(db.Bolt.Stats())
	// var stats_json string
	// stats = json.Unmarshal(stats, &stats_json)
	f.Println("Here are the stats -> ", stats, err)

}

// func addPatient(db *storm.DB, p Patient) error {

// }

// func getPatient(db *storm.DB, patient_key []byte, p Patient) error {

// }

func main() {

	models.InitDB("hc_db_init.json")
	deleteFile("hc_db_init.json")

	//start_server()

	// r := chi.NewRouter()
	// r.Use(middleware.RequestID)
	// r.Use(middleware.Logger)
	// r.Use(middleware.Logger)
	// r.Use(middleware.Recoverer)

	// r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Welcome to HealthChain :)\n"))
	// })

	// r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("pong"))
	// })

	// r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
	// 	panic("test")
	// })

	// db, err := storm.Open("hc.db")
	// if err != nil {
	// 	fmt.Println("db.Open exception")
	// }
	// defer db.Close()

	// view_db(db)

	// init_db()

	// var test Patient
	// get_err := db.One("PatientKey", "Laurie Feliciano6361781291938-11-18US", &test)

	// if get_err != nil {
	// 	fmt.Print(get_err.Error())
	// }
	// fmt.Println("This should be a record :) -> ", test)
	// http.ListenAndServe(":3333", r)

}

func deleteFile(filepath string) {
	var err = os.Remove(filepath)
	if isError(err) {
		return
	}

}

func isError(err error) bool {
	if err != nil {
		f.Println(err.Error())
	}

	return (err != nil)
}
