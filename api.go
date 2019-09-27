package main


import (
	"net/http"
	log "github.com/sirupsen/logrus"
	"github.com/go-chi/chi"

)

func api(){

  r := chi.NewRouter()

  err := http.ListenAndServe(":3000", r)
  if err != nil{
  
    log.WithError(err).Error("unable to listen and serve")
  }

}
