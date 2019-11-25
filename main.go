package main

// main
// Runs discovery procedure to find peers in network and starts webserver
func main() {

	// use the same database object for the peer driver
	pd := CreatePeerDriver()
	go pd.Discovery()

	// as well as the api
	api := NewAPI(pd.uuid)
	api.Run()

	return
}
