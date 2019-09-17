package main

import ( 
	"fmt"
	"flag"
	"time"
	// "os"
	// "net"
	"log"
	"github.com/grandcat/zeroconf"
)


var (
	name     = flag.String("name", "GoZeroconfGo", "The name for the service.")
	service  = flag.String("service", "_workstation._tcp", "Set the service type of the new service.")
	domain   = flag.String("domain", "local.", "Set the network domain. Default should be fine.")
	port     = flag.Int("port", 42424, "Set the port the service is listening to.")
	waitTime = flag.Int("wait", 10, "Duration in [s] to publish service for.")
)

func main(){
	flag.Parse()

	fmt.Println("Starting discovery")

	// Discover all services on the network (e.g. _workstation._tcp)
	// resolver, err := zeroconf.NewResolver(nil)
	// if err != nil {
	// 		log.Fatalln("Failed to initialize resolver:", err.Error())
	// }

	// register service
	server, err := zeroconf.Register("HealthChain", "_healthchain._tcp", "local.", 42424, []string{"txtv=0", "lo=1", "la=2"}, nil)
	if err != nil {
			panic(err)
	}

	defer server.Shutdown()
	log.Println("Published service:")
	log.Println("- Name:", *name)
	log.Println("- Type:", *service)
	log.Println("- Domain:", *domain)
	log.Println("- Port:", *port)

	// fmt.Println(resolver)
	time.Sleep(5 * time.Second)

	return 
}