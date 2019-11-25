package main

// main
// Runs discovery procedure to find peers in network and starts webserver
func main() {

	// start peer discovery in background
	pd := CreatePeerDriver()
	go pd.Discovery()

	// create API as the main thread service
	api := NewAPI(pd.uuid)
	api.Run()

	return
}
