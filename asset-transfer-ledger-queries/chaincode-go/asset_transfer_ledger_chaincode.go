/*
 SPDX-License-Identifier: Apache-2.0
*/

/*
====CHAINCODE EXECUTION SAMPLES (CLI) ==================

==== Invoke assets ====
peer chaincode invoke -C myc1 -n asset_transfer -c '{"Args":["CreateAsset","asset1","blue","5","tom","35"]}'
peer chaincode invoke -C myc1 -n asset_transfer -c '{"Args":["CreateAsset","asset2","red","4","tom","50"]}'
peer chaincode invoke -C myc1 -n asset_transfer -c '{"Args":["CreateAsset","asset3","blue","6","tom","70"]}'
peer chaincode invoke -C myc1 -n asset_transfer -c '{"Args":["TransferAsset","asset2","jerry"]}'
peer chaincode invoke -C myc1 -n asset_transfer -c '{"Args":["TransferAssetByColor","blue","jerry"]}'
peer chaincode invoke -C myc1 -n asset_transfer -c '{"Args":["DeleteAsset","asset1"]}'

==== Query assets ====
peer chaincode query -C myc1 -n asset_transfer -c '{"Args":["ReadAsset","asset1"]}'
peer chaincode query -C myc1 -n asset_transfer -c '{"Args":["GetAssetsByRange","asset1","asset3"]}'
peer chaincode query -C myc1 -n asset_transfer -c '{"Args":["GetAssetHistory","asset1"]}'

Rich Query (Only supported if CouchDB is used as state database):
peer chaincode query -C myc1 -n asset_transfer -c '{"Args":["QueryAssetsByOwner","tom"]}'
peer chaincode query -C myc1 -n asset_transfer -c '{"Args":["QueryAssets","{\"selector\":{\"owner\":\"tom\"}}"]}'

Rich Query with Pagination (Only supported if CouchDB is used as state database):
peer chaincode query -C myc1 -n asset_transfer -c '{"Args":["QueryAssetsWithPagination","{\"selector\":{\"owner\":\"tom\"}}","3",""]}'

INDEXES TO SUPPORT COUCHDB RICH QUERIES

Indexes in CouchDB are required in order to make JSON queries efficient and are required for
any JSON query with a sort. Indexes may be packaged alongside
chaincode in a META-INF/statedb/couchdb/indexes directory. Each index must be defined in its own
text file with extension *.json with the index definition formatted in JSON following the
CouchDB index JSON syntax as documented at:
http://docs.couchdb.org/en/2.3.1/api/database/find.html#db-index

This asset transfer ledger example chaincode demonstrates a packaged
index which you can find in META-INF/statedb/couchdb/indexes/indexOwner.json.

If you have access to the your peer's CouchDB state database in a development environment,
you may want to iteratively test various indexes in support of your chaincode queries.  You
can use the CouchDB Fauxton interface or a command line curl utility to create and update
indexes. Then once you finalize an index, include the index definition alongside your
chaincode in the META-INF/statedb/couchdb/indexes directory, for packaging and deployment
to managed environments.

In the examples below you can find index definitions that support asset transfer ledger
chaincode queries, along with the syntax that you can use in development environments
to create the indexes in the CouchDB Fauxton interface or a curl command line utility.


Index for docType, owner.

Example curl command line to define index in the CouchDB channel_chaincode database
curl -i -X POST -H "Content-Type: application/json" -d "{\"index\":{\"fields\":[\"docType\",\"owner\"]},\"name\":\"indexOwner\",\"ddoc\":\"indexOwnerDoc\",\"type\":\"json\"}" http://hostname:port/myc1_assets/_index


Index for docType, owner, size (descending order).

Example curl command line to define index in the CouchDB channel_chaincode database:
curl -i -X POST -H "Content-Type: application/json" -d "{\"index\":{\"fields\":[{\"size\":\"desc\"},{\"docType\":\"desc\"},{\"owner\":\"desc\"}]},\"ddoc\":\"indexSizeSortDoc\", \"name\":\"indexSizeSortDesc\",\"type\":\"json\"}" http://hostname:port/myc1_assets/_index

Rich Query with index design doc and index name specified (Only supported if CouchDB is used as state database):
peer chaincode query -C myc1 -n asset_transfer -c '{"Args":["QueryAssets","{\"selector\":{\"docType\":\"asset\",\"owner\":\"tom\"}, \"use_index\":[\"_design/indexOwnerDoc\", \"indexOwner\"]}"]}'

Rich Query with index design doc specified only (Only supported if CouchDB is used as state database):
peer chaincode query -C myc1 -n asset_transfer -c '{"Args":["QueryAssets","{\"selector\":{\"docType\":{\"$eq\":\"asset\"},\"owner\":{\"$eq\":\"tom\"},\"size\":{\"$gt\":0}},\"fields\":[\"docType\",\"owner\",\"size\"],\"sort\":[{\"size\":\"desc\"}],\"use_index\":\"_design/indexSizeSortDoc\"}"]}'
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const index = "color~name"
const EblIndex = "originCompanyID~eblNo"

// SimpleChaincode implements the fabric-contract-api-go programming model
type SimpleChaincode struct {
	contractapi.Contract
}

type Asset struct {
	DocType        string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	ID             string `json:"ID"`      //the field tags are needed to keep case from bouncing around
	Color          string `json:"color"`
	Size           int    `json:"size"`
	Owner          string `json:"owner"`
	AppraisedValue int    `json:"appraisedValue"`
}

type Ebl struct {
	EblNo                  string  `json:"eblNo"`
	OriginCompanyID        string  `json:"originCompanyID"`
	OriginCompanyName      string  `json:"originCompanyName"`
	ShipperCompanyID       string  `json:"shipperCompanyID"`
	ShipperCompanyName     string  `json:"shipperCompanyName"`
	ConsigneeCompanyID     string  `json:"consigneeCompanyID"`
	ConsigneeCompanyName   string  `json:"consigneeCompanyName"`
	NotifyPartyCompanyID   string  `json:"notifyPartyCompanyID"`
	NotifyPartyCompanyName string  `json:"notifyPartyCompanyName"`
	PlaceOfReceipt         string  `json:"placeOfReceipt"`
	OceanVessel            string  `json:"oceanVessel"`
	PortOfLoading          string  `json:"portOfLoading"`
	PortOfDescharge        string  `json:"portOfDescharge"`
	PlaceOfDestination     string  `json:"placeOfDestination"`
	PlaceOfDelivery        string  `json:"placeOfDelivery"`
	ShippingMarkes         string  `json:"shippingMarkes"`
	QuantityOfPackages     float64 `json:"quantityOfPackages"`
	KindOfPackagesGW       string  `json:"kindOfPackagesGw"`
	KindOfPackagesM        string  `json:"kindOfPackagesM"`
	DescriptionOfGoods     string  `json:"descriptionOfGoods"`
	GrossWeight            float64 `json:"grossWeight"`
	Measurement            float64 `json:"measurement"`
	FreightAndCharges      string  `json:"freightAndCharges"`
	PlaceOfIssue           string  `json:"placeOfIssue"`
	DateOfIssue            int64   `json:"dateOfIssue"`
	DeliveryAgent          string  `json:"deliveryAgent"`
	ShippedOnBoard         int64   `json:"shippedOnBoard"`
	NumOfEBL               int64   `json:"numOfEbl"`
	DateOfIssueDeadline    int64   `json:"dateOfIssueDeadline"`
	Status                 string  `json:"status"`
	File                   string  `json:"file"`
	ContractFiles          string  `json:"contractFiles"`
	InvoiceFiles           string  `json:"invoiceFiles"`
	TransferCompanyID      string  `json:"transferCompanyID"`
	TransferCompanyName    string  `json:"transferCompanyName"`
	CompanyID              int64   `json:"companyID"`
	CompanyName            string  `json:"companyName"`
}

// HistoryQueryResult structure used for returning result of history query
type HistoryQueryResult struct {
	Record    *Asset    `json:"record"`
	TxId      string    `json:"txId"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"isDelete"`
}

// PaginatedQueryResult structure used for returning paginated query results and metadata
type PaginatedQueryResult struct {
	Records             []*Asset `json:"records"`
	FetchedRecordsCount int32    `json:"fetchedRecordsCount"`
	Bookmark            string   `json:"bookmark"`
}

type EblHistoryQueryResult struct {
	Record    *Ebl      `json:"record"`
	TxId      string    `json:"txId"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"isDelete"`
}

type EblPaginatedQueryResult struct {
	Records             []*Ebl `json:"records"`
	FetchedRecordsCount int32  `json:"fetchedRecordsCount"`
	Bookmark            string `json:"bookmark"`
}

// CreateAsset initializes a new asset in the ledger
func (t *SimpleChaincode) CreateAsset(ctx contractapi.TransactionContextInterface, assetID, color string, size int, owner string, appraisedValue int) error {
	exists, err := t.AssetExists(ctx, assetID)
	if err != nil {
		return fmt.Errorf("failed to get asset: %v", err)
	}
	if exists {
		return fmt.Errorf("asset already exists: %s", assetID)
	}

	asset := &Asset{
		DocType:        "asset",
		ID:             assetID,
		Color:          color,
		Size:           size,
		Owner:          owner,
		AppraisedValue: appraisedValue,
	}
	assetBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(assetID, assetBytes)
	if err != nil {
		return err
	}

	//  Create an index to enable color-based range queries, e.g. return all blue assets.
	//  An 'index' is a normal key-value entry in the ledger.
	//  The key is a composite key, with the elements that you want to range query on listed first.
	//  In our case, the composite key is based on indexName~color~name.
	//  This will enable very efficient state range queries based on composite keys matching indexName~color~*
	colorNameIndexKey, err := ctx.GetStub().CreateCompositeKey(index, []string{asset.Color, asset.ID})
	if err != nil {
		return err
	}
	//  Save index entry to world state. Only the key name is needed, no need to store a duplicate copy of the asset.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	return ctx.GetStub().PutState(colorNameIndexKey, value)
}

// CreateEbl initializes a new EBL (Electronic Bill of Lading) in the ledger
func (t *SimpleChaincode) CreateEbl(ctx contractapi.TransactionContextInterface, eblNo, originCompanyID, originCompanyName, shipperCompanyID, shipperCompanyName,
	consigneeCompanyID, consigneeCompanyName, notifyPartyCompanyID, notifyPartyCompanyName, placeOfReceipt, oceanVessel, portOfLoading, portOfDescharge, placeOfDestination,
	placeOfDelivery, shippingMarkes, contractFiles, invoiceFiles, transferCompanyID, transferCompanyName, kindOfPackagesGW, kindOfPackagesM, descriptionOfGoods, deliveryAgent,
	companyName, freightAndCharges, status, file, placeOfIssue string,
	quantityOfPackages, grossWeight, measurement float64,
	dateOfIssue, shippedOnBoard, numOfEBL, dateOfIssueDeadline, companyID int64) error {
	// Check if the EBL already exists
	exists, err := t.EblExists(ctx, eblNo)
	if err != nil {
		return fmt.Errorf("failed to check if EBL exists: %v", err)
	}
	if exists {
		return fmt.Errorf("EBL already exists: %s", eblNo)
	}

	// Create the EBL object
	ebl := &Ebl{
		EblNo:                  eblNo,
		OriginCompanyID:        originCompanyID,
		OriginCompanyName:      originCompanyName,
		ShipperCompanyID:       shipperCompanyID,
		ShipperCompanyName:     shipperCompanyName,
		ConsigneeCompanyID:     consigneeCompanyID,
		ConsigneeCompanyName:   consigneeCompanyName,
		NotifyPartyCompanyID:   notifyPartyCompanyID,
		NotifyPartyCompanyName: notifyPartyCompanyName,
		PlaceOfReceipt:         placeOfReceipt,
		OceanVessel:            oceanVessel,
		PortOfLoading:          portOfLoading,
		PortOfDescharge:        portOfDescharge,
		PlaceOfDestination:     placeOfDestination,
		PlaceOfDelivery:        placeOfDelivery,
		ShippingMarkes:         shippingMarkes,
		QuantityOfPackages:     quantityOfPackages,
		KindOfPackagesGW:       kindOfPackagesGW,
		KindOfPackagesM:        kindOfPackagesM,
		DescriptionOfGoods:     descriptionOfGoods,
		GrossWeight:            grossWeight,
		Measurement:            measurement,
		FreightAndCharges:      freightAndCharges,
		PlaceOfIssue:           placeOfIssue,
		DateOfIssue:            dateOfIssue,
		DeliveryAgent:          deliveryAgent,
		ShippedOnBoard:         shippedOnBoard,
		NumOfEBL:               numOfEBL,
		DateOfIssueDeadline:    dateOfIssueDeadline,
		Status:                 "Created",
		File:                   file,
		ContractFiles:          contractFiles,
		InvoiceFiles:           invoiceFiles,
		TransferCompanyID:      transferCompanyID,
		TransferCompanyName:    transferCompanyName,
		CompanyName:            companyName,
		CompanyID:              companyID,
	}

	// Marshal the EBL struct into JSON bytes
	eblBytes, err := json.Marshal(ebl)
	if err != nil {
		return err
	}

	// Save the EBL to the world state
	err = ctx.GetStub().PutState(eblNo, eblBytes)
	if err != nil {
		return err
	}

	// Create an index for EBL by EblNo to support efficient querying based on EBL number
	// Here we create an index for querying by EblNo
	eblNoIndexKey, err := ctx.GetStub().CreateCompositeKey(EblIndex, []string{ebl.EblNo, ebl.OriginCompanyID})
	if err != nil {
		return err
	}

	// Save the index entry to the world state
	value := []byte{0x00} // we don't need to store a duplicate copy of the EBL, just the index key
	err = ctx.GetStub().PutState(eblNoIndexKey, value)
	if err != nil {
		return err
	}

	// Optionally, create other indexes based on different fields (e.g. status, companyID, etc.)

	return nil
}

func (t *SimpleChaincode) SubmitEbl(ctx contractapi.TransactionContextInterface, eblNo string) error {
	ebl, err := t.ReadEbl(ctx, eblNo)
	if err != nil {
		return err
	}

	// Update the EBL status to "Submitted"
	ebl.Status = "Submitted"

	// Marshal the EBL struct into JSON bytes
	eblBytes, err := json.Marshal(ebl)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(eblNo, eblBytes)
}

// ReadAsset retrieves an asset from the ledger
func (t *SimpleChaincode) ReadAsset(ctx contractapi.TransactionContextInterface, assetID string) (*Asset, error) {
	assetBytes, err := ctx.GetStub().GetState(assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset %s: %v", assetID, err)
	}
	if assetBytes == nil {
		return nil, fmt.Errorf("asset %s does not exist", assetID)
	}

	var asset Asset
	err = json.Unmarshal(assetBytes, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// ReadEbl retrieves an EBL from the ledger based on originCompanyID.
func (t *SimpleChaincode) ReadEbl(ctx contractapi.TransactionContextInterface, eblNo string) (*Ebl, error) {
	// Get the EBL from the world state
	eblBytes, err := ctx.GetStub().GetState(eblNo)
	if err != nil {
		return nil, fmt.Errorf("failed to read EBL %s from world state: %v", eblNo, err)
	}
	if eblBytes == nil {
		return nil, fmt.Errorf("EBL %s does not exist", eblNo)
	}

	// Unmarshal the EBL bytes into an EBL struct
	var ebl Ebl
	err = json.Unmarshal(eblBytes, &ebl)
	if err != nil {
		return nil, err
	}

	return &ebl, nil
}

// DeleteAsset removes an asset key-value pair from the ledger
func (t *SimpleChaincode) DeleteAsset(ctx contractapi.TransactionContextInterface, assetID string) error {
	asset, err := t.ReadAsset(ctx, assetID)
	if err != nil {
		return err
	}

	err = ctx.GetStub().DelState(assetID)
	if err != nil {
		return fmt.Errorf("failed to delete asset %s: %v", assetID, err)
	}

	colorNameIndexKey, err := ctx.GetStub().CreateCompositeKey(index, []string{asset.Color, asset.ID})
	if err != nil {
		return err
	}

	// Delete index entry
	return ctx.GetStub().DelState(colorNameIndexKey)
}

// TransferAsset transfers an asset by setting a new owner name on the asset
func (t *SimpleChaincode) TransferAsset(ctx contractapi.TransactionContextInterface, assetID, newOwner string) error {
	asset, err := t.ReadAsset(ctx, assetID)
	if err != nil {
		return err
	}

	asset.Owner = newOwner
	assetBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(assetID, assetBytes)
}

// constructQueryResponseFromIterator constructs a slice of assets from the resultsIterator
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*Asset, error) {
	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		fmt.Println("Query Result:", string(queryResult.Value)) // 打印返回的原始值
		if len(queryResult.Value) == 0 {
			return nil, fmt.Errorf("empty value found for asset ID: %s", queryResult.Key)
		}
		var asset Asset
		err = json.Unmarshal(queryResult.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

func constructEblQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*Ebl, error) {
	var ebls []*Ebl
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var ebl Ebl
		err = json.Unmarshal(queryResult.Value, &ebl)
		if err != nil {
			return nil, err
		}
		ebls = append(ebls, &ebl)
	}

	return ebls, nil
}

func (t *SimpleChaincode) GetAssetsByRange(ctx contractapi.TransactionContextInterface, startKey, endKey string) ([]*Asset, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}

func (t *SimpleChaincode) TransferAssetByColor(ctx contractapi.TransactionContextInterface, color, newOwner string) error {
	// Execute a key range query on all keys starting with 'color'
	coloredAssetResultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(index, []string{color})
	if err != nil {
		return err
	}
	defer coloredAssetResultsIterator.Close()

	for coloredAssetResultsIterator.HasNext() {
		responseRange, err := coloredAssetResultsIterator.Next()
		if err != nil {
			return err
		}

		_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(responseRange.Key)
		if err != nil {
			return err
		}

		if len(compositeKeyParts) > 1 {
			returnedAssetID := compositeKeyParts[1]
			asset, err := t.ReadAsset(ctx, returnedAssetID)
			if err != nil {
				return err
			}
			asset.Owner = newOwner
			assetBytes, err := json.Marshal(asset)
			if err != nil {
				return err
			}
			err = ctx.GetStub().PutState(returnedAssetID, assetBytes)
			if err != nil {
				return fmt.Errorf("transfer failed for asset %s: %v", returnedAssetID, err)
			}
		}
	}

	return nil
}

func (t *SimpleChaincode) QueryAssetsByOwner(ctx contractapi.TransactionContextInterface, owner string) ([]*Asset, error) {
	queryString := fmt.Sprintf(`{"selector":{"docType":"asset","owner":"%s"}}`, owner)
	return getQueryResultForQueryString(ctx, queryString)
}

func (t *SimpleChaincode) QueryAssets(ctx contractapi.TransactionContextInterface, queryString string) ([]*Asset, error) {
	return getQueryResultForQueryString(ctx, queryString)
}

func getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Asset, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}

func (t *SimpleChaincode) GetAssetsByRangeWithPagination(ctx contractapi.TransactionContextInterface, startKey string, endKey string, pageSize int, bookmark string) (*PaginatedQueryResult, error) {

	resultsIterator, responseMetadata, err := ctx.GetStub().GetStateByRangeWithPagination(startKey, endKey, int32(pageSize), bookmark)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	assets, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}
	return &PaginatedQueryResult{
		Records:             assets,
		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}

func (t *SimpleChaincode) GetEblByRangeWithPagination(ctx contractapi.TransactionContextInterface, startKey string, endKey string, pageSize int, bookmark string) (*EblPaginatedQueryResult, error) {
	resultsIterator, responseMetadata, err := ctx.GetStub().GetStateByRangeWithPagination(startKey, endKey, int32(pageSize), bookmark)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	ebls, err := constructEblQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}
	return &EblPaginatedQueryResult{
		Records:             ebls,
		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}

func (t *SimpleChaincode) QueryAssetsWithPagination(ctx contractapi.TransactionContextInterface, queryString string, pageSize int, bookmark string) (*PaginatedQueryResult, error) {

	return getQueryResultForQueryStringWithPagination(ctx, queryString, int32(pageSize), bookmark)
}

func (t *SimpleChaincode) QueryEblWithPagination(ctx contractapi.TransactionContextInterface, queryString string, pageSize int, bookmark string) (*EblPaginatedQueryResult, error) {
	return getEblQueryResultForQueryStringWithPagination(ctx, queryString, int32(pageSize), bookmark)
}

func getQueryResultForQueryStringWithPagination(ctx contractapi.TransactionContextInterface, queryString string, pageSize int32, bookmark string) (*PaginatedQueryResult, error) {

	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	assets, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	return &PaginatedQueryResult{
		Records:             assets,
		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}

func getEblQueryResultForQueryStringWithPagination(ctx contractapi.TransactionContextInterface, queryString string, pageSize int32, bookmark string) (*EblPaginatedQueryResult, error) {
	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	ebls, err := constructEblQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}
	if ebls == nil {
		return nil, fmt.Errorf("No EBLs found")
	}
	return &EblPaginatedQueryResult{
		Records:             ebls,
		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}

// GetAssetHistory returns the chain of custody for an asset since issuance.
func (t *SimpleChaincode) GetAssetHistory(ctx contractapi.TransactionContextInterface, assetID string) ([]HistoryQueryResult, error) {
	log.Printf("GetAssetHistory: ID %v", assetID)

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(assetID)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &asset)
			if err != nil {
				return nil, err
			}
		} else {
			asset = Asset{
				ID: assetID,
			}
		}

		timestamp, err := ptypes.Timestamp(response.Timestamp)
		if err != nil {
			return nil, err
		}

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: timestamp,
			Record:    &asset,
			IsDelete:  response.IsDelete,
		}
		records = append(records, record)
	}

	return records, nil
}

// AssetExists returns true when asset with given ID exists in the ledger.
func (t *SimpleChaincode) AssetExists(ctx contractapi.TransactionContextInterface, assetID string) (bool, error) {
	assetBytes, err := ctx.GetStub().GetState(assetID)
	if err != nil {
		return false, fmt.Errorf("failed to read asset %s from world state. %v", assetID, err)
	}

	return assetBytes != nil, nil
}

// EblExists returns true when EBL with given EblNo exists in the ledger.
func (t *SimpleChaincode) EblExists(ctx contractapi.TransactionContextInterface, eblNo string) (bool, error) {
	eblBytes, err := ctx.GetStub().GetState(eblNo)
	if err != nil {
		return false, fmt.Errorf("failed to read EBL %s from world state. %v", eblNo, err)
	}

	return eblBytes != nil, nil
}

// InitLedger creates the initial set of assets in the ledger.
func (t *SimpleChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{
		{DocType: "asset", ID: "asset1", Color: "blue", Size: 5, Owner: "Tomoko", AppraisedValue: 300},
		{DocType: "asset", ID: "asset2", Color: "red", Size: 5, Owner: "Brad", AppraisedValue: 400},
		{DocType: "asset", ID: "asset3", Color: "green", Size: 10, Owner: "Jin Soo", AppraisedValue: 500},
		{DocType: "asset", ID: "asset4", Color: "yellow", Size: 10, Owner: "Max", AppraisedValue: 600},
		{DocType: "asset", ID: "asset5", Color: "black", Size: 15, Owner: "Adriana", AppraisedValue: 700},
		{DocType: "asset", ID: "asset6", Color: "white", Size: 15, Owner: "Michel", AppraisedValue: 800},
	}

	for _, asset := range assets {
		err := t.CreateAsset(ctx, asset.ID, asset.Color, asset.Size, asset.Owner, asset.AppraisedValue)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SimpleChaincode{})
	if err != nil {
		log.Panicf("Error creating asset chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting asset chaincode: %v", err)
	}
}
