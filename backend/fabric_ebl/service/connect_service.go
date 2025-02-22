package service

import (
	"context"
	"fabric_ebl/domain"
	"fabric_ebl/repo"
	"fabric_ebl/sal/jwt"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"go.uber.org/dig"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type ConnectService interface {
	Contract(ctx context.Context, userName string, companyName string) (contract *gateway.Contract, closeFunc CloseFunc, err error)
	ParseToken(ctx context.Context, req string) (token string, claims map[string]any, userId int64, user *domain.UserDO, company *domain.CompanyDO, err error)
}

type ConnectServiceParam struct {
	dig.In
	FabricEblRepo repo.FabricEblRepo
}

type ConnectServiceImpl struct {
	p ConnectServiceParam
}

func (c ConnectServiceImpl) ParseToken(ctx context.Context, reqToken string) (respToken string, claims map[string]any, userId int64, user *domain.UserDO, company *domain.CompanyDO, err error) {
	token := reqToken
	claims, err = jwt.ValidateToken(ctx, token)
	if err != nil {
		logger.CtxErrorf(ctx, "ParseToken failed, err = %v", err)
		return "", nil, 0, nil, nil, err
	}
	userId, err = strconv.ParseInt(claims["user_id"].(string), 10, 64)
	if err != nil {
		logger.CtxErrorf(ctx, "ParseInt failed, err = %v", err)
		return "", nil, 0, nil, nil, err
	}
	user, err = c.p.FabricEblRepo.QueryUserById(ctx, userId)
	if err != nil {
		logger.CtxErrorf(ctx, "QueryUserById failed, err = %v", err)
		return "", nil, 0, nil, nil, err
	}
	company, err = c.p.FabricEblRepo.QueryCompanyById(ctx, user.CompanyID)
	if err != nil {
		logger.CtxErrorf(ctx, "QueryCompanyById failed, err = %v", err)
		return "", nil, 0, nil, nil, err
	}
	return token, claims, userId, user, company, nil
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
