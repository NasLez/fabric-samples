package service

import (
	"context"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"go.uber.org/dig"
	"log"
	"os"
	"path/filepath"
)

type ConnectService interface {
	Contract(ctx context.Context, userName string, companyName string) (contract *gateway.Contract, closeFunc CloseFunc, err error)
}

type ConnectServiceParam struct {
	dig.In
}

type ConnectServiceImpl struct {
	p ConnectServiceParam
}

func (c ConnectServiceImpl) Contract(ctx context.Context, userName string, companyName string) (contract *gateway.Contract, closeFunc CloseFunc, err error) {
	err = os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Printf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v\n", err)
		return nil, nil, err
	}
	walletName := companyName + "_" + userName
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Printf("Failed to create wallet: %v\n", err)
		return nil, nil, err
	}
	if !wallet.Exists(walletName) {
		err = addUserToWallet(wallet, walletName)
		if err != nil {
			log.Printf("Failed to populate wallet contents: %v\n", err)
			return nil, nil, err
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
		gateway.WithIdentity(wallet, walletName),
	)
	if err != nil {
		log.Printf("Failed to connect to gateway: %v\n", err)
		return nil, nil, err
	}
	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		log.Printf("Failed to get network: %v\n", err)
		gw.Close()
		return nil, nil, err
	}
	contract = network.GetContract("basic")
	return contract, gw.Close, nil
}

type CloseFunc func()

func NewConnectService(p ConnectServiceParam) ConnectService {
	return &ConnectServiceImpl{
		p: p,
	}
}
