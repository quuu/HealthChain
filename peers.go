package main

import (
	// "fmt"
	"flag"
	"time"

	// "os"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/grandcat/zeroconf"
	log "github.com/sirupsen/logrus"
)

var (
	name     = flag.String("name", "GoZeroconfGo", "The name for the service.")
	service  = flag.String("service", "_workstation._tcp", "Set the service type of the new service.")
	domain   = flag.String("domain", "local.", "Set the network domain. Default should be fine.")
	port     = flag.Int("port", 42424, "Set the port the service is listening to.")
	waitTime = flag.Int("wait", 10, "Duration in [s] to publish service for.")
)

func discovery() {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.WithError(err).Error("unable to create listener")
		return
	}
	port := listener.Addr().(*net.TCPAddr).Port
	log.Printf("Listening at %s", listener.Addr())

	r := chi.NewRouter()
	//start the server in the background
	go func() {
		err := http.Serve(listener, r)
		log.WithError(err).Error("unable to listen and serve")
	}()

	// register service
	server, err := zeroconf.Register("HealthChain", "_healthchain._tcp", "local.", port, []string{"txtv=0", "lo=1", "la=2"}, nil)
	if err != nil {
		panic(err)
	}

	defer server.Shutdown()

	// Clean exit.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	select {
	case <-sig:
		// Exit by user
	case <-time.After(time.Second * 120):
		// Exit by timeout
	}

	log.Println("Shutting down.")

}
