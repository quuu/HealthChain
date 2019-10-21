package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

type Peer struct {
	Addresses []net.IP
	Port      int
	ID        string
	Expires   time.Time
}

type DiscoveryManager struct {
	m     *sync.Mutex
	peers map[string]*Peer
}

func main() {
	dm := &DiscoveryManager{
		m:     &sync.Mutex{},
		peers: map[string]*Peer{},
	}

	r := chi.NewRouter()
	r.Post("/notify", dm.notifyHandler)
	r.Get("/fetch", dm.fetchHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}

func (dm *DiscoveryManager) fetchHandler(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	err := enc.Encode(dm.peers)
	if err != nil {

		log.WithError(err).Error("unable to encode")
		http.Error(w, "unable to encode", http.StatusInternalServerError)
	}

}

func (dm *DiscoveryManager) notifyHandler(w http.ResponseWriter, r *http.Request) {
	peer := &Peer{}

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(peer)
	if err != nil {
		log.WithError(err).Error("unable to decode")
		http.Error(w, "unable to decode", http.StatusInternalServerError)
		return
	}

	peer.Expires = time.Now().Add(time.Second * 30)
	dm.m.Lock()
	dm.peers[peer.ID] = peer
	peers := []*Peer{}
	for peerID, peer := range dm.peers {
		if !peer.Expires.After(time.Now()) {
			delete(dm.peers, peerID)
			continue
		}
		peers = append(peers, peer)
	}

	dm.m.Unlock()

	enc := json.NewEncoder(w)
	err = enc.Encode(peers)
	if err != nil {
		log.WithError(err).Error("unable to encode")
		http.Error(w, "unable to encode", http.StatusInternalServerError)
	}
}
