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

// Peer struct used to store information about individual
type Peer struct {
	Addresses []net.IP
	Port      int
	ID        string
	Expires   time.Time
}

// DiscoveryDriver struct is used
type DiscoveryDriver struct {
	m     *sync.Mutex
	peers map[string]*Peer
}

func main() {
	dm := &DiscoveryDriver{
		m:     &sync.Mutex{},
		peers: map[string]*Peer{},
	}

	r := chi.NewRouter()

	// route responsible for saving a new client's ip addresses
	r.Post("/notify", dm.notifyHandler)

	// route responsible for debugging and seeing the active clients in the network
	r.Get("/fetch", dm.fetchHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}

func (dm *DiscoveryDriver) fetchHandler(w http.ResponseWriter, r *http.Request) {

	// simply decode the peers and push it to http responsewriter
	enc := json.NewEncoder(w)
	err := enc.Encode(dm.peers)
	if err != nil {

		log.WithError(err).Error("unable to encode")
		http.Error(w, "unable to encode", http.StatusInternalServerError)
	}

}

func (dm *DiscoveryDriver) notifyHandler(w http.ResponseWriter, r *http.Request) {
	peer := &Peer{}

	// set up getting the peer information to save to the driver
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(peer)
	if err != nil {
		log.WithError(err).Error("unable to decode")
		http.Error(w, "unable to decode", http.StatusInternalServerError)
		return
	}

	// set the expiration for 30 seconds if they don't communicate back
	peer.Expires = time.Now().Add(time.Second * 30)
	dm.m.Lock()
	dm.peers[peer.ID] = peer
	peers := []*Peer{}

	// check for dead peers, clear out the storage
	for peerID, peer := range dm.peers {
		if !peer.Expires.After(time.Now()) {
			delete(dm.peers, peerID)
			continue
		}
		peers = append(peers, peer)
	}

	dm.m.Unlock()

	// write to the endpoint the list of all clients
	enc := json.NewEncoder(w)
	err = enc.Encode(peers)
	if err != nil {
		log.WithError(err).Error("unable to encode")
		http.Error(w, "unable to encode", http.StatusInternalServerError)
	}
}
