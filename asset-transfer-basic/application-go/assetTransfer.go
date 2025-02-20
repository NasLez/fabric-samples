/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func main() {
	log.Println("============ application-golang starts ============")

	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}

	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}

	if !wallet.Exists("nas1") {
		err = populateWallet(wallet)
		if err != nil {
			log.Fatalf("Failed to populate wallet contents: %v", err)
		}
	}

	ccpPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"connection-org1.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "nas1"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}

	contract := network.GetContract("basic")

	//log.Println("--> Submit Transaction: InitLedger, function creates the initial set of assets on the ledger")
	//result, err := contract.SubmitTransaction("InitLedger")
	//if err != nil {
	//	log.Fatalf("Failed to Submit transaction: %v", err)
	//}
	//log.Println(string(result))

	//log.Println("--> Evaluate Transaction: GetAllAssets, function returns all the current assets on the ledger")
	//result, err = contract.EvaluateTransaction("GetAllAssets")
	//if err != nil {
	//	log.Fatalf("Failed to evaluate transaction: %v", err)
	//}
	//log.Println(string(result))

	log.Println("--> Submit Transaction: CreateAsset, creates new asset with ID, color, owner, size, and appraisedValue arguments")
	result, err := contract.SubmitTransaction("CreateAsset", "asset13", "yellow", "5", "Tom", "1300")
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}
	log.Println(string(result))

	log.Println("--> Evaluate Transaction: ReadAsset, function returns an asset with a given assetID")
	result, err = contract.EvaluateTransaction("ReadAsset", "asset13")
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v\n", err)
	}
	log.Println(string(result))

	log.Println("--> Evaluate Transaction: AssetExists, function returns 'true' if an asset with given assetID exist")
	result, err = contract.EvaluateTransaction("AssetExists", "asset1")
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v\n", err)
	}
	log.Println(string(result))

	log.Println("--> Submit Transaction: TransferAsset asset1, transfer to new owner of Tom")
	_, err = contract.SubmitTransaction("TransferAsset", "asset1", "Tom")
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}

	log.Println("--> Evaluate Transaction: ReadAsset, function returns 'asset1' attributes")
	result, err = contract.EvaluateTransaction("ReadAsset", "asset1")
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}
	log.Println(string(result))

	log.Println("--> Submit Transaction: CreateEbl, creates new EBL with provided details")
	result, err = contract.SubmitTransaction(
		"CreateEbl",               // chaincode method
		"ebl123",                  // eblNo
		"company1",                // originCompanyID
		"Company A",               // originCompanyName
		"company2",                // shipperCompanyID
		"Company B",               // shipperCompanyName
		"company3",                // consigneeCompanyID
		"Company C",               // consigneeCompanyName
		"company4",                // notifyPartyCompanyID
		"Company D",               // notifyPartyCompanyName
		"NYC",                     // placeOfReceipt
		"Vessel X",                // oceanVessel
		"Port A",                  // portOfLoading
		"Port B",                  // portOfDescharge
		"Place A",                 // placeOfDestination
		"Place B",                 // placeOfDelivery
		"Mark A",                  // shippingMarkes
		"contractfile1",           // contractFiles (Array of files)
		"invoicefile1",            // invoiceFiles (Array of files)
		"company5",                // transferCompanyID
		"Company E",               // transferCompanyName
		"Package A",               // kindOfPackagesGW
		"Kg",                      // kindOfPackagesM
		"Sample Description",      // descriptionOfGoods
		"Agent A",                 // deliveryAgent
		"Company X",               // companyName
		"Freight charges details", // freightAndCharges
		"Active",                  // status
		"filehash123",             // file
		"Place X",                 // placeOfIssue
		"100.5",                   // quantityOfPackages
		"200.0",                   // grossWeight
		"300.0",                   // measurement
		"1635734400",              // dateOfIssue (int64 timestamp)
		"1635734400",              // shippedOnBoard (int64 timestamp)
		"1",                       // numOfEBL
		"1635734400",              // dateOfIssueDeadline (int64 timestamp)
		"1234567",                 // companyID (int64)
	)
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}
	log.Println(string(result))

	// 查询交易：读取 EBL
	log.Println("--> Evaluate Transaction: ReadEbl, function returns EBL with given originCompanyID")
	result, err = contract.EvaluateTransaction("ReadEbl", "ebl123")
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v\n", err)
	}
	log.Println(string(result))
	log.Println("============ application-golang ends ============")
}

func populateWallet(wallet *gateway.Wallet) error {
	log.Println("============ Populating wallet ============")
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
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	return wallet.Put("nas1", identity)
}
