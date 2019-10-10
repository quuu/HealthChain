package handlers

import (
	f "fmt"
	"net/http"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	name := r.URL.Query().Get("name")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(f.Sprintf("Hello %s!", name)))
}
