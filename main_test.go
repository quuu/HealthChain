package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/grandcat/zeroconf"
)

func TestGetHash(t *testing.T) {
	// hash := GetHash("Testing", "Account", "USA", "123")
	// resultHash := []byte("071cb9c3296810ec371341b0df272044c31c8fec9fb1617ea2d8ae5a56681317")

	// if string(hash) != string(resultHash) {

	// 	t.Errorf("Expected %s\n Got %s\n", resultHash, hash)
	// }

}
func TestEncryption(t *testing.T) {

	hash := GetHash("first", "last", "country", "code")

	message := []byte("this is a record")

	encrypted_message := Encrypt(hash, message)

	if string(encrypted_message) == string(message) {
		t.Errorf("Encryption did not alter the message\n")
	}
}

func TestDecryption(t *testing.T) {

	hash := GetHash("first", "last", "country", "code")

	message := []byte("this is a record")

	encrypted_message := Encrypt(hash, message)

	decrypted_message := Decrypt(hash, encrypted_message)

	if string(decrypted_message) != string(message) {
		t.Errorf("Encryption and decryption are not symmetric\n")
	}
}

func TestDatabaseOperations(t *testing.T) {

	db, err := storm.Open("test.db")
	if err != nil {
		t.Errorf("Could not create database in running directory\n")
	}

	type Data struct {
		ID      string `storm:"id"`
		Message string
	}

	toStore := &Data{ID: "1", Message: "First record!"}
	err = db.Save(toStore)
	if err != nil {
		t.Errorf("Could not store into database object\n")
	}

	err = os.Remove("./test.db")
	if err != nil {
		t.Errorf("Could not remove test databse in running directory\n")
	}

}

func TestPeerDiscovery(t *testing.T) {
	server, err := zeroconf.Register("TestingZeroconf", "_healthchain._tcp", "local.", 4000, nil, nil)

	if err != nil {

		t.Errorf("Failed to register service\n")
	}
	defer server.Shutdown()

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		t.Errorf("Failed to create resolver\n")
	}
	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	err = resolver.Browse(ctx, "_workstation._tcp", "local.", entries)
	if err != nil {
		t.Errorf("Failed to browse\n")
	}

	<-ctx.Done()
}

func TestServiceRegistration(t *testing.T) {

	server, err := zeroconf.Register("TestingZeroconf", "_healthchain._tcp", "local.", 4000, nil, nil)

	if err != nil {

		t.Errorf("failed to register service")
	}
	defer server.Shutdown()

}

func TestAPI(t *testing.T) {

}
