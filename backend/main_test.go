package main

import (
	"bytes"
	"context"
	"github.com/asdine/storm/v3"
	"github.com/grandcat/zeroconf"
	"os"
	"testing"
	"time"
)

// TestGetHash
// Proves GetHash is deterministic
func TestGetHash(t *testing.T) {
	hash_1 := GetHash("Testing", "Account", "USA", "123")
	hash_2 := GetHash("Testing", "Account", "USA", "123")
	if !bytes.Equal(hash_1, hash_2) {
		t.Errorf("Expected %s\n Got %s\n", string(hash_2), string(hash_1))
	}

}

// TestEncryption
// Proves that encryption process makes data unreadable
func TestEncryption(t *testing.T) {

	hash := GetHash("first", "last", "country", "code")

	message := []byte("this is a record")

	encrypted_message := Encrypt(hash, message)

	if string(encrypted_message) == string(message) {
		t.Errorf("Encryption did not alter the message\n")
	}
}

// TestDecryption
// Proves that Decryption returns original data before encryption
func TestDecryption(t *testing.T) {

	hash := GetHash("first", "last", "country", "code")

	message := []byte("this is a record")

	encrypted_message := Encrypt(hash, message)

	decrypted_message := Decrypt(hash, encrypted_message)

	if string(decrypted_message) != string(message) {
		t.Errorf("Encryption and decryption are not symmetric\n")
	}
}

//  TestDatabaseOperations
// Proves that data can be saves into db succesfully
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

// TestPeerDiscovery
// Proves that peer discovery finds peers succesfully
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

// TestServiceRegistration
// Proves that  Zeroconf registeres correctly to tcp protocol
func TestServiceRegistration(t *testing.T) {

	server, err := zeroconf.Register("TestingZeroconf", "_healthchain._tcp", "local.", 4000, nil, nil)
	defer server.Shutdown()
	if err != nil {
		t.Errorf("failed to register service")
	}

}
