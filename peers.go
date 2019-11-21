package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"strconv"

	"github.com/go-chi/chi"
	"github.com/grandcat/zeroconf"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

// PeerDriver is a struct used to manage concurrency and a list of peers
type PeerDriver struct {
	uuid  string
	m     *sync.Mutex
	me    *Peer
	peers map[string]*Peer
	api   *API
}

// Peer is a helper struct to store information about the peer
type Peer struct {
	ID        string
	Addresses []net.IP
	Port      int
}

// function that handles displaying all the messages
// TODO - NOTE -> this method should be fine for get/set
// integrate with the storage
func (pd *PeerDriver) recordHandler(w http.ResponseWriter, r *http.Request) {

	// host all the
	w.Header().Set("Content-Type", "application/json")
	var records []*EncryptedRecord
	db := PublicDB()
	defer db.Close()
	err := db.All(&records)
	fmt.Println(records)
	b, err := json.MarshalIndent(records, "", " ")
	if err != nil {
		panic(err.Error())
	}

	w.Write(b)
}

func (pd *PeerDriver) peerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// create a new encoder
	enc := json.NewEncoder(w)

	// make sure no one else is editing the peers while this happens
	pd.m.Lock()
	defer pd.m.Unlock()

	// encode the peers to the response writer
	if err := enc.Encode(pd.peers); err != nil {
		http.Error(w, "unable to encode peers", http.StatusInternalServerError)
		log.WithError(err).Error("unable to encode peers")
		return
	}

}

// Create a PeerDriver object
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
	r.Get("/records", pd.recordHandler)
	r.Get("/peers", pd.peerHandler)

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

	// handle global entries from the web
	globalEntries := make(chan *Peer)
	go func() {
		c := http.Client{
			Timeout: 5 * time.Second,
		}
		ticker := time.Tick(time.Second * 5)
		for range ticker {
			pd.m.Lock()
			/*
				if pd.me == nil {
					continue
				}

			*/
			b, err := json.Marshal(pd.me)
			pd.m.Unlock()
			if err != nil {
				log.WithError(err).Error("unable to marshal")
				continue
			}
			resp, err := c.Post("https://trans-dogfish-253118.appspot.com/notify", "application/json", bytes.NewBuffer(b))
			if err != nil {
				log.WithError(err).Error("unable to post")
				continue
			}

			peers := []*Peer{}
			dec := json.NewDecoder(resp.Body)
			err = dec.Decode(&peers)
			if err != nil {
				log.WithError(err).Error("unable to decode peers")
				continue
			}
			for _, peer := range peers {
				log.Println(peer.ID + " " + strconv.Itoa(peer.Port))
				globalEntries <- peer
			}
		}
	}()

	// inifinte loop waiting for more entries
	ticker := time.Tick(1 * time.Second)
	for {
		select {
		case entry := <-entries:
			pd.handleEntry(entry)
		case <-ticker:
			pd.fetchRecords()
		case entry := <-globalEntries:
			pd.handleGlobalEntry(entry)
		}
	}
}

func (pd *PeerDriver) handleGlobalEntry(entry *Peer) {

	// ignore if self
	if entry.ID == pd.uuid {
		return
	}

	// create a new peer
	p := &Peer{
		ID:        entry.ID,
		Port:      entry.Port,
		Addresses: entry.Addresses,
	}

	pd.m.Lock()
	defer pd.m.Unlock()
	if peer, ok := pd.peers[p.ID]; ok {
		log.Printf("peer %s already known", p.ID)
		peer.Addresses = p.Addresses
		return
	}

	log.Printf("new peer: %+v", p)
	pd.peers[p.ID] = p
}

// function responsible for receiving a new peer
// adding it to the list of peers
func (pd *PeerDriver) handleEntry(entry *zeroconf.ServiceEntry) {

	// if found self, add addresses if not already in the list
	if entry.Instance == pd.uuid {
		pd.m.Lock()
		pd.me = &Peer{
			ID:   entry.Instance,
			Port: entry.Port,
		}

		// more addresses for self
		pd.me.Addresses = append(pd.me.Addresses, entry.AddrIPv6...)
		pd.me.Addresses = append(pd.me.Addresses, entry.AddrIPv4...)

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

	c := http.Client{
		Timeout: time.Second * 5,
	}

	pd.m.Lock()
	defer pd.m.Unlock()

	// for all known peers
	for index, peer := range pd.peers {

		retrieved := false

		// for every address the peer has
		for _, addr := range peer.Addresses {

			endpoint := url.URL{
				Scheme: "http",
				Path:   "/records",
			}

			// check what kind of address is being used
			if addr.To4() != nil {
				endpoint.Host = addr.String() + ":" + strconv.Itoa(peer.Port)
			} else if addr.To16() != nil {
				endpoint.Host = "[" + addr.String() + "]:" + strconv.Itoa(peer.Port)
			} else {
				log.Errorf("peer %s addr type is unknown", index)
				continue
			}

			// create the request
			req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
			if err != nil {
				continue
			}

			// actually do the request
			resp, err := c.Do(req)
			if err != nil {
				continue
			}

			// got something
			retrieved = true

			// parse the data
			dec := json.NewDecoder(resp.Body)

			// turn the response into encrypted record objects
			var encrypted_records []*EncryptedRecord
			err = dec.Decode(&encrypted_records)
			if err != nil {
				panic(err.Error())
			}

			// call function used to store the records that are unique
			pd.handleRecords(encrypted_records)

			log.Println("got records")

		}
		// if nothing found, connection is dead
		if !retrieved {
			delete(pd.peers, index)
			log.Printf("peer was dead %s", index)
		}

	}
}

// function that will loop through an array of records
// 	and save each into storm
func (pd *PeerDriver) handleRecords(encrypted_records []*EncryptedRecord) {

	// db := GetDB()
	// defer db.Close()
	var records_temp []Record
	for _, rec := range encrypted_records {

		p := &Patient{}
		p = GetPatient(string(rec.PatientID))

		// if the patient doesn't exist
		if p == nil {
			log.Println("Patient does not exist, making new patient")
			p := Patient{PatientKey: string(rec.PatientID), Records: records_temp, Node: "hc_1"}
			AddPatient(p)
		} else {

			found := true
			for _, record := range p.Records {
				if found {

				}
				fmt.Println("record is " + string(record.Message))
				fmt.Println("encrypted record is " + string(rec.Contents))
				// found = false
				// if(record ==)
			}
			// otherwise, check the records and make sure it's not a duplicate

		}

		// err := db.Save(rec)
		// if err != nil {
		// 	panic(err.Error())
		// }
	}
}
