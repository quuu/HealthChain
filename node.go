package main

import (
	// utils
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	// "bytes"
	// "encoding/json"

	//chi-chi
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	//skv
)

// notes on progress
// establish nosql way to store data

// rest
// post -> add document to db assigned with hash value
// get -> get all documents associated with hash value

func main() {

	file, _ := os.Open("./data/KaggleV2-May-2016.csv")
	fmt.Print(file)
	// store, _ := skv.Open("./local_hc.db")

	// store.Put("hash1", file)
	// store.Get("hash1", &file)
	// p

}
func start_server() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "files")
	FileServer(r, "/files", http.Dir(filesDir))
	fmt.Println("Node up on port 3333")
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
