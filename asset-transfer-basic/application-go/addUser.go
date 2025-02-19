package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

//
//func main() {
//	log.Println("============ application-golang starts ============")
//
//	// Define a command line flag for username
//	username := flag.String("username", "", "The username to populate the wallet for")
//	flag.Parse()
//
//	// Ensure that the username is provided
//	if *username == "" {
//		log.Fatal("Username is required. Use the -username flag to provide it.")
//	}
//
//	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
//	if err != nil {
//		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environment variable: %v", err)
//	}
//
//	wallet, err := gateway.NewFileSystemWallet("wallet")
//	if err != nil {
//		log.Fatalf("Failed to create wallet: %v", err)
//	}
//
//	// Check if the wallet contains the identity for the provided username
//	if !wallet.Exists(*username) {
//		err = addUserToWallet(wallet, *username)
//		if err != nil {
//			log.Fatalf("Failed to populate wallet contents: %v", err)
//		}
//	}
//
//	log.Println("============ application-golang ends ============")
//}

func addUserToWallet(wallet *gateway.Wallet, username string) error {
	log.Println("============ Populating wallet for user:", username, "===========")
	credPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	// Read the certificate PEM
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// There's a single file in this directory containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have exactly one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	// Store the identity in the wallet under the provided username
	return wallet.Put(username, identity)
}
