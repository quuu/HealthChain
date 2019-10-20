package main

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/grandcat/zeroconf"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

// PeerDriver is a struct used to manage concurrency and a list of peers
type PeerDriver struct {
	uuid  string
	m     *sync.Mutex
	self  *Peer
	peers map[string]*Peer
}

// Peer is a helper struct to store information about the peer
type Peer struct {
	ID        string
	Addresses []net.IP
	Port      int
}

// function that handles displaying all the messages
// TODO
// integrate with the storage
func recordHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Records go here"))
}

func peerHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("peers go here"))
}

func CreatePeerDriver() *PeerDriver {
	pd := &PeerDriver{
		m:     &sync.Mutex{},
		peers: map[string]*Peer{},
	}

	return pd

}

func (pd *PeerDriver) Discovery() {

	// initialize all the endpoints to serve publicly
	r := chi.NewRouter()
	r.Get("/records", recordHandler)
	r.Get("/peers", peerHandler)

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.WithError(err).Error("unable to create listener")
		return
	}

	// save the port for registering zeroconf
	port := listener.Addr().(*net.TCPAddr).Port

	// serve the public record
	log.Debugf("PeerManager listening at %s", listener.Addr())
	go func() {
		err := http.Serve(listener, r)
		log.WithError(err).Error("PeerManager unable to listen and serve")
	}()

	// create a unique identifier for this node
	u := uuid.NewV4()

	// save it
	pd.uuid = u.String()

	// register it with the unique name
	go func() {
		for {
			server, err := zeroconf.Register(u.String(), "_healthchain._tcp", "local.", port, nil, nil)
			if err != nil {
				log.Fatal(err)
				continue
			}
			<-time.After(time.Second * 5)
			server.Shutdown()
		}
	}()

	log.Printf("started listening at %d", port)

	// now browse for other services

	// will store the new peers discovered
	entries := make(chan *zeroconf.ServiceEntry)
	resolver, nil := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	// get the running go routine for registration
	ctx := context.Background()
	err = resolver.Browse(ctx, "_healthchain._tcp", "local.", entries)
	if err != nil {
		log.WithError(err).Error("Unable to browse")
		return
	}

	// HANDLE GLOBAL ENTRIES WITH CLOUD HERE
	// TODO

	// inifinte loop waiting for more entries
	ticker := time.Tick(1 * time.Second)
	for {
		select {
		case entry := <-entries:
			pd.handleEntry(entry)
		case <-ticker:
			pd.fetchRecords()
		}
	}
}

// function responsible for receiving a new peer
// adding it to the list of peers
func (pd *PeerDriver) handleEntry(entry *zeroconf.ServiceEntry) {

	// if foudn self, add addresses if not already in the list
	if entry.Instance == pd.uuid {
		pd.m.Lock()
		pd.self = &Peer{
			ID:   entry.Instance,
			Port: entry.Port,
		}
		pd.self.Addresses = append(pd.self.Addresses, entry.AddrIPv6...)
		pd.self.Addresses = append(pd.self.Addresses, entry.AddrIPv4...)

		pd.m.Unlock()
		log.Println("found self")
		return
	}

	//got a unique entry
	log.Println("got an entry")
	p := &Peer{
		ID:   entry.Instance,
		Port: entry.Port,
	}
	p.Addresses = append(p.Addresses, entry.AddrIPv6...)
	p.Addresses = append(p.Addresses, entry.AddrIPv4...)

	pd.m.Lock()
	defer pd.m.Unlock()

	// already discovered, just add more addresses
	if peer, ok := pd.peers[p.ID]; ok {

		log.Printf("peer %s already discovered", p.ID)
		peer.Addresses = p.Addresses
		return
	}

	// found completely new one
	log.Printf("found new peer %+v", p)
	pd.peers[p.ID] = p

}

// function responsible for asking peers for records
// TODO
// only fetch on a need to use basis
// currently set to fetch every 1 second when theres no new entries
func (pd *PeerDriver) fetchRecords() {

	log.Println("currently fetching records ")
}
