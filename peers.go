package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	//	"strconv"
	"time"
	//	"time"

	// "os"
	//	"net"
	//	"os"
	//	"os/signal"
	//	"syscall"

	"github.com/go-chi/chi"
	"github.com/grandcat/zeroconf"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

func recordHandler(w http.ResponseWriter, r *http.Request) {

}

func discovery() {

	r := chi.NewRouter()
	r.Get("/records", recordHandler)

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

	// register it with the unique name
	service, err := zeroconf.Register(u.String(), "_healthchain._tcp", "local.", port, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer service.Shutdown()

	// now browse for other services
	resolver, nil := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatal(err)
	}

	// make a channel to save the results
	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			fmt.Println(entry)
			fmt.Println(entry.Port)
		}
		fmt.Println("out of entries")
	}(entries)

	// get the background process to browse from
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(10))
	defer cancel()
	err = resolver.Browse(ctx, "_healthchain._tcp", "local.", entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
	}
	<-ctx.Done()

	// HANDLE GLOBAL ENTRIES WITH CLOUD HERE
	// TODO
	ticker := time.Tick(1 * time.Second)
	for {
		select {
		case entry := <-entries:
			handleEntry(entry)
		case <-ticker:
			fetchMessages()
		}
	}
}

func handleEntry(entry *zeroconf.ServiceEntry) {

}

func fetchMessages() {

}
