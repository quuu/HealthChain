package main

import (
	"context"
	"flag"
	"fmt"
	"time"
	//	"time"

	// "os"
	//	"net"
	//	"net/http"
	//	"os"
	//	"os/signal"
	//	"syscall"

	//	"github.com/go-chi/chi"
	"github.com/grandcat/zeroconf"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

var (
	name     = flag.String("name", "GoZeroconfGo", "The name for the service.")
	service  = flag.String("service", "_workstation._tcp", "Set the service type of the new service.")
	domain   = flag.String("domain", "local.", "Set the network domain. Default should be fine.")
	port     = flag.Int("port", 42424, "Set the port the service is listening to.")
	waitTime = flag.Int("wait", 10, "Duration in [s] to publish service for.")
)

func main() {

	// create a unique identifier for this node
	u := uuid.NewV4()

	// register it with the unique name
	service, err := zeroconf.Register(u.String(), "_healthchain._tcp", "local.", 8080, nil, nil)
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
		}
		fmt.Println("out of entries")
	}(entries)

	// get the background process to browse from
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(*waitTime))
	defer cancel()
	err = resolver.Browse(ctx, "_healthchain._tcp", "local.", entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
	}

	<-ctx.Done()
	// Wait some additional time to see debug messages on go routine shutdown.
	time.Sleep(1 * time.Second)

	select {}
}
