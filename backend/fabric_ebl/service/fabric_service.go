package service

import (
	"context"
	"encoding/json"
	"fabric_ebl/biz_error"
	"fabric_ebl/common/id_gen"
	"fabric_ebl/domain/converter"
	"fabric_ebl/repo"
	"fabric_ebl/sal/jwt"
	"fabric_ebl/sal/rpc/fabric_ipfs_rpc"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/wxl-server/idl_gen/kitex_gen/fabric_ebl"
	"github.com/wxl-server/idl_gen/kitex_gen/fabric_ipfs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/wxl-server/idl_gen/kitex_gen/common_user"
	"go.uber.org/dig"
)

type FabricEblService interface {
	UpdatePassword(ctx context.Context, req *common_user.UpdatePasswordReq) (*common_user.UpdatePasswordResp, error)
	Login(ctx context.Context, req *fabric_ebl.LoginReq) (resp *fabric_ebl.LoginResp, err error)
	CreateCompany(ctx context.Context, req *fabric_ebl.CreateCompanyReq) (resp *fabric_ebl.CreateCompanyResp, err error)
	GetUserInfo(ctx context.Context, req *fabric_ebl.GetUserInfoReq) (*fabric_ebl.GetUserInfoResp, error)
	GetCompanyAllList(ctx context.Context, req *fabric_ebl.GetCompanyAllListReq) (*fabric_ebl.GetCompanyAllListResp, error)
	CreateEbl(ctx context.Context, req *fabric_ebl.CreateEblReq) (*fabric_ebl.CreateEblResp, error)
	QueryAllEblList(ctx context.Context, req *fabric_ebl.QueryAllEblListReq) (*fabric_ebl.QueryAllEblListResp, error)
	QueryEblList(ctx context.Context, req *fabric_ebl.QueryEblListReq) (*fabric_ebl.QueryEblListResp, error)
}

type Param struct {
	dig.In
	FabricEblRepo repo.FabricEblRepo
}

type FabricEblServiceImpl struct {
	p Param
}

func NewUserService(p Param) FabricEblService {
	return &FabricEblServiceImpl{
		p: p,
	}
}

func (u FabricEblServiceImpl) QueryEblList(ctx context.Context, req *fabric_ebl.QueryEblListReq) (*fabric_ebl.QueryEblListResp, error) {
	token := req.Token
	claims, err := jwt.ValidateToken(ctx, token)
	if err != nil {
		logger.CtxErrorf(ctx, "ParseToken failed, err = %v", err)
		return nil, biz_error.ParseTokenError
	}
	userId, err := strconv.ParseInt(claims["user_id"].(string), 10, 64)
	user, err := u.p.FabricEblRepo.QueryUserById(ctx, userId)
	if err != nil {
		logger.CtxErrorf(ctx, "QueryUserById failed, err = %v", err)
		return nil, err
	}
	company, err := u.p.FabricEblRepo.QueryCompanyById(ctx, user.CompanyID)
	if err != nil {
		logger.CtxErrorf(ctx, "QueryCompanyById failed, err = %v", err)
		return nil, err
	}
	err = os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Println("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}
	walletName := company.Name + "_" + user.Name
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Println("Failed to create wallet: %v", err)
	}
	if !wallet.Exists(walletName) {
		err = addUserToWallet(wallet, walletName)
		if err != nil {
			log.Println("Failed to populate wallet contents: %v", err)
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
		log.Println("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()
	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		log.Println("Failed to get network: %v", err)
	}
	contract := network.GetContract("basic")
	log.Println("--> Submit Transaction: GetEblByRangeWithPagination, creates new EBL with provided details")

	contractPageSize := strconv.FormatInt(*req.PageSize, 10)
	selector := generateSelectorString(req)
	//selector = "{\"selector\":{\"portOfDescharge\":\"Port B\"},\"use_index\":[\"_design/indexEblDoc\",\"indexEbl\"]}"
	log.Println(selector)
	result, err := contract.SubmitTransaction("QueryEblWithPagination", selector, contractPageSize, *req.Bookmark)
	if err != nil {
		log.Println("Failed to Submit transaction: %v", err)
		return nil, nil
	}
	log.Println(string(result))
	resp, err := GetEblByRangeWithPaginationResp2DO(result)
	if err != nil {
		log.Println("Failed to convert result: %v", err)
	}
	return &fabric_ebl.QueryEblListResp{
		EblList:             resp.EblList,
		Bookmark:            resp.Bookmark,
		FetchedRecordsCount: resp.FetchedRecordsCount,
	}, nil
}

func (u FabricEblServiceImpl) QueryAllEblList(ctx context.Context, req *fabric_ebl.QueryAllEblListReq) (*fabric_ebl.QueryAllEblListResp, error) {
	token := req.Token
	claims, err := jwt.ValidateToken(ctx, token)
	if err != nil {
		logger.CtxErrorf(ctx, "ParseToken failed, err = %v", err)
		return nil, biz_error.ParseTokenError
	}
	userId, err := strconv.ParseInt(claims["user_id"].(string), 10, 64)
	user, err := u.p.FabricEblRepo.QueryUserById(ctx, userId)
	if err != nil {
		logger.CtxErrorf(ctx, "QueryUserById failed, err = %v", err)
		return nil, err
	}
	company, err := u.p.FabricEblRepo.QueryCompanyById(ctx, user.CompanyID)
	if err != nil {
		logger.CtxErrorf(ctx, "QueryCompanyById failed, err = %v", err)
		return nil, err
	}
	err = os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Println("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}
	walletName := company.Name + "_" + user.Name
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Println("Failed to create wallet: %v", err)
	}
	if !wallet.Exists(walletName) {
		err = addUserToWallet(wallet, walletName)
		if err != nil {
			log.Println("Failed to populate wallet contents: %v", err)
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
		log.Println("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()
	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		log.Println("Failed to get network: %v", err)
	}
	contract := network.GetContract("basic")
	log.Println("--> Submit Transaction: GetEblByRangeWithPagination, creates new EBL with provided details")

	contractPageSize := strconv.FormatInt(*req.PageSize, 10)
	result, err := contract.SubmitTransaction("GetEblByRangeWithPagination", "", "", contractPageSize, *req.Bookmark)
	if err != nil {
		log.Println("Failed to Submit transaction: %v", err)
	}
	log.Println(string(result))
	resp, err := GetEblByRangeWithPaginationResp2DO(result)
	if err != nil {
		log.Println("Failed to convert result: %v", err)
	}
	return resp, nil
}

func (u FabricEblServiceImpl) CreateEbl(ctx context.Context, req *fabric_ebl.CreateEblReq) (*fabric_ebl.CreateEblResp, error) {
	token := req.Token
	claims, err := jwt.ValidateToken(ctx, token)
	if err != nil {
		logger.CtxErrorf(ctx, "ParseToken failed, err = %v", err)
		return nil, biz_error.ParseTokenError
	}
	userId, err := strconv.ParseInt(claims["user_id"].(string), 10, 64)
	user, err := u.p.FabricEblRepo.QueryUserById(ctx, userId)
	if err != nil {
		logger.CtxErrorf(ctx, "QueryUserById failed, err = %v", err)
		return nil, err
	}
	company, err := u.p.FabricEblRepo.QueryCompanyById(ctx, user.CompanyID)
	if err != nil {
		logger.CtxErrorf(ctx, "QueryCompanyById failed, err = %v", err)
		return nil, err
	}

	{
		req.Ebl.OriginCompanyID = strconv.FormatInt(company.ID, 10)
		req.Ebl.OriginCompanyName = company.Name
		req.Ebl.CompanyID = company.ID
		req.Ebl.CompanyName = company.Name
		req.Ebl.Status = "created"
	}

	{
		reqShipperCompanyID, err := strconv.ParseInt(req.Ebl.ShipperCompanyID, 10, 64)
		shipperCompany, err := u.p.FabricEblRepo.QueryCompanyById(ctx, reqShipperCompanyID)
		if err != nil {
			logger.CtxErrorf(ctx, "QueryCompanyById failed, err = %v", err)
			return nil, err
		}
		req.Ebl.ShipperCompanyName = shipperCompany.Name
	}

	{
		reqConsigneeCompanyID, err := strconv.ParseInt(req.Ebl.ConsigneeCompanyID, 10, 64)
		consigneeCompany, err := u.p.FabricEblRepo.QueryCompanyById(ctx, reqConsigneeCompanyID)
		if err != nil {
			logger.CtxErrorf(ctx, "QueryCompanyById failed, err = %v", err)
			return nil, err
		}
		req.Ebl.ConsigneeCompanyName = consigneeCompany.Name
	}

	{
		reqNotifyPartyCompanyID, err := strconv.ParseInt(req.Ebl.NotifyPartyCompanyID, 10, 64)
		notifyPartyCompany, err := u.p.FabricEblRepo.QueryCompanyById(ctx, reqNotifyPartyCompanyID)
		if err != nil {
			logger.CtxErrorf(ctx, "QueryCompanyById failed, err = %v", err)
			return nil, err
		}
		req.Ebl.NotifyPartyCompanyName = notifyPartyCompany.Name
	}

	err = os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Println("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}
	walletName := company.Name + "_" + user.Name
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Println("Failed to create wallet: %v", err)
	}
	if !wallet.Exists(walletName) {
		err = addUserToWallet(wallet, walletName)
		if err != nil {
			log.Println("Failed to populate wallet contents: %v", err)
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
		log.Println("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()
	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		log.Println("Failed to get network: %v", err)
	}
	contract := network.GetContract("basic")
	log.Println("--> Submit Transaction: CreateEbl, creates new EBL with provided details")
	ID, err := id_gen.NextID()
	req.Ebl.EblNo = strconv.FormatInt(ID, 10)
	{
		resp, err := fabric_ipfs_rpc.CreateEblDocx(ctx, reqHttp2Rpc(req))
		if err != nil {
			logger.CtxErrorf(ctx, "CreateEblDocx failed, err = %v", err)
			return nil, err
		}
		logger.Infof("CreateEblDocx success, resp = %v", resp.FileHash)
		req.Ebl.File = resp.FileHash
	}
	result, err := contract.SubmitTransaction(
		"CreateEbl",                              // chaincode method
		req.Ebl.EblNo,                            // eblNo
		req.Ebl.OriginCompanyID,                  // originCompanyID
		req.Ebl.OriginCompanyName,                // originCompanyName
		req.Ebl.ShipperCompanyID,                 // shipperCompanyID
		req.Ebl.ShipperCompanyName,               // shipperCompanyName
		req.Ebl.ConsigneeCompanyID,               // consigneeCompanyID
		req.Ebl.ConsigneeCompanyName,             // consigneeCompanyName
		req.Ebl.NotifyPartyCompanyID,             // notifyPartyCompanyID
		req.Ebl.NotifyPartyCompanyName,           // notifyPartyCompanyName
		req.Ebl.PlaceOfReceipt,                   // placeOfReceipt
		req.Ebl.OceanVessel,                      // oceanVessel
		req.Ebl.PortOfLoading,                    // portOfLoading
		req.Ebl.PortOfDescharge,                  // portOfDescharge
		req.Ebl.PlaceOfDestination,               // placeOfDestination
		req.Ebl.PlaceOfDelivery,                  // placeOfDelivery
		req.Ebl.ShippingMarkes,                   // shippingMarkes
		strings.Join(req.Ebl.ContractFiles, ";"), // contractFiles (can be a file or file path)
		strings.Join(req.Ebl.InvoiceFiles, ";"),  // invoiceFiles (can be a file or file path)
		"",                                       // transferCompanyID
		"",                                       // transferCompanyName
		req.Ebl.KindOfPackagesGW,                 // kindOfPackagesGW
		req.Ebl.KindOfPackagesM,                  // kindOfPackagesM
		req.Ebl.DescriptionOfGoods,               // descriptionOfGoods
		req.Ebl.DeliveryAgent,                    // deliveryAgent
		req.Ebl.CompanyName,                      // companyName
		req.Ebl.FreightAndCharges,                // freightAndCharges
		req.Ebl.Status,                           // status
		req.Ebl.File,                             // file
		req.Ebl.PlaceOfIssue,                     // placeOfIssue
		strconv.FormatFloat(req.Ebl.QuantityOfPackages, 'f', -1, 64), // quantityOfPackages
		strconv.FormatFloat(req.Ebl.GrossWeight, 'f', -1, 64),        // grossWeight
		strconv.FormatFloat(req.Ebl.Measurement, 'f', -1, 64),        // measurement
		strconv.FormatInt(req.Ebl.DateOfIssue, 10),                   // dateOfIssue
		strconv.FormatInt(req.Ebl.ShippedOnBoard, 10),                // shippedOnBoard
		strconv.FormatInt(req.Ebl.NumOfEBL, 10),                      // numOfEBL
		strconv.FormatInt(req.Ebl.DateOfIssueDeadline, 10),           // dateOfIssueDeadline
		strconv.FormatInt(req.Ebl.CompanyID, 10),                     // companyID
	)
	if err != nil {
		log.Println("Failed to Submit transaction: %v", err)
	}
	log.Println(string(result))

	// 查询交易：读取 EBL
	log.Println("--> Evaluate Transaction: ReadEbl, function returns EBL with given eblNo")
	result, err = contract.EvaluateTransaction("ReadEbl", req.Ebl.EblNo)
	if err != nil {
		log.Println("Failed to evaluate transaction: %v\n", err)
	}
	log.Println(string(result))
	return &fabric_ebl.CreateEblResp{
		Id: ID,
	}, nil
}

func (u FabricEblServiceImpl) GetCompanyAllList(ctx context.Context, req *fabric_ebl.GetCompanyAllListReq) (*fabric_ebl.GetCompanyAllListResp, error) {
	companyList, err := u.p.FabricEblRepo.QueryCompanyAll(ctx)
	if err != nil {
		logger.CtxErrorf(ctx, "QueryCompanyAll failed, err = %v", err)
		return nil, err
	}
	var companyListResp []*fabric_ebl.Company
	for _, company := range companyList {
		companyListResp = append(companyListResp, &fabric_ebl.Company{
			Id:          company.ID,
			CompanyCode: company.Code,
			CompanyName: company.Name,
			CompanyType: fabric_ebl.CompanyType(company.Type),
		})
	}
	return &fabric_ebl.GetCompanyAllListResp{
		CompanyList: companyListResp,
	}, nil
}

func (u FabricEblServiceImpl) GetUserInfo(ctx context.Context, req *fabric_ebl.GetUserInfoReq) (*fabric_ebl.GetUserInfoResp, error) {
	token := req.Token
	claims, err := jwt.ValidateToken(ctx, token)
	if err != nil {
		logger.CtxErrorf(ctx, "ParseToken failed, err = %v", err)
		return nil, biz_error.ParseTokenError
	}
	user_id, err := strconv.ParseInt(claims["user_id"].(string), 10, 64)
	user, err := u.p.FabricEblRepo.QueryUserById(ctx, user_id)
	if err != nil {
		logger.CtxErrorf(ctx, "QueryUserById failed, err = %v", err)
		return nil, err
	}
	company, err := u.p.FabricEblRepo.QueryCompanyById(ctx, user.CompanyID)
	if err != nil {
		logger.CtxErrorf(ctx, "QueryCompanyById failed, err = %v", err)
		return nil, err
	}
	return &fabric_ebl.GetUserInfoResp{
		UserId:      user.ID,
		UserName:    user.Name,
		UserEmail:   user.Email,
		UserType:    fabric_ebl.UserType(user.Type),
		CompanyId:   company.ID,
		CompanyName: company.Name,
		CompanyCode: company.Code,
		CompanyType: fabric_ebl.CompanyType(company.Type),
	}, nil
}

func (u FabricEblServiceImpl) SignUp(ctx context.Context, req *common_user.SignUpReq) (resp *common_user.SignUpResp, err error) {
	count, err := u.p.FabricEblRepo.CountUser(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		logger.CtxErrorf(ctx, "sign up failed, email has been used, email = %s", req.Email)
		return nil, biz_error.SignUpError
	}
	id, err := u.p.FabricEblRepo.CreateUser(ctx, converter.SignUpReqDTO2DO(req))
	if err != nil {
		logger.CtxErrorf(ctx, "CreateUser failed, err = %v", err)
		return nil, err
	}
	return &common_user.SignUpResp{
		Id: id,
	}, nil
}

func (u FabricEblServiceImpl) UpdatePassword(ctx context.Context, req *common_user.UpdatePasswordReq) (*common_user.UpdatePasswordResp, error) {
	if req.OldPassword == req.Password {
		logger.CtxErrorf(ctx, "UpdatePassword failed, err = %v", biz_error.UpdatePasswordError)
		return nil, biz_error.UpdatePasswordError
	}
	do, err := u.p.FabricEblRepo.QueryUser(ctx, req.Email)
	if err != nil {
		logger.CtxErrorf(ctx, "QueryUser failed, err = %v", err)
		return nil, err
	}
	if do.Password != req.OldPassword {
		logger.CtxErrorf(ctx, "UpdatePassword failed, err = %v", biz_error.UpdatePasswordError2)
		return nil, biz_error.UpdatePasswordError2
	}
	do.Password = req.Password
	err = u.p.FabricEblRepo.UpdatePassword(ctx, do)
	if err != nil {
		logger.CtxErrorf(ctx, "UpdatePassword failed, err = %v", err)
		return nil, err
	}
	return &common_user.UpdatePasswordResp{}, nil
}

func (u FabricEblServiceImpl) Login(ctx context.Context, req *fabric_ebl.LoginReq) (resp *fabric_ebl.LoginResp, err error) {
	do, err := u.p.FabricEblRepo.QueryUser(ctx, req.Email)
	if err != nil {
		logger.CtxErrorf(ctx, "QueryUser failed, err = %v", err)
		return nil, biz_error.LoginError
	}
	if do.Password != req.Password {
		logger.CtxErrorf(ctx, "login failed, password is wrong")
		return nil, biz_error.LoginError2
	}
	token, err := jwt.GenerateToken(ctx, do.ID)
	if err != nil {
		logger.CtxErrorf(ctx, "login failed, err = %v", err)
		return nil, biz_error.LoginError
	}
	return &fabric_ebl.LoginResp{
		Token: token,
	}, nil
}

func (u FabricEblServiceImpl) CreateCompany(ctx context.Context, req *fabric_ebl.CreateCompanyReq) (resp *fabric_ebl.CreateCompanyResp, err error) {
	{
		count, err := u.p.FabricEblRepo.CountCompanyByCode(ctx, req.CompanyCode)
		if err != nil {
			logger.CtxErrorf(ctx, "QueryCompanyByCode failed, err = %v", err)
			return nil, err
		}
		if count > 0 {
			logger.CtxErrorf(ctx, "CreateCompany failed, company code has been created, companyCode = %s", req.CompanyCode)
			return nil, biz_error.CreateCompanyError1
		}
	}
	{
		count, err := u.p.FabricEblRepo.CountUser(ctx, req.AdminEmail)
		if err != nil {
			logger.CtxErrorf(ctx, "QueryUser failed, err = %v", err)
			return nil, err
		}
		if count > 0 {
			logger.CtxErrorf(ctx, "CreateCompany failed, user has been created, email = %s", req.AdminEmail)
			return nil, biz_error.CreateCompanyError2
		}
	}
	companyId, err := u.p.FabricEblRepo.CreateCompany(ctx, converter.CreateCompanyReqDTO2DO(req))
	if err != nil {
		logger.CtxErrorf(ctx, "CreateCompany failed, err = %v", err)
		return nil, err
	}
	userId, err := u.p.FabricEblRepo.CreateUser(ctx, converter.CreateCompanyReqDTO2UserDO(req, companyId))
	if err != nil {
		logger.CtxErrorf(ctx, "CreateUser failed, err = %v", err)
		return nil, err
	}
	{
		wallet, err := gateway.NewFileSystemWallet("wallet")
		if err != nil {
			log.Println("Failed to create wallet: %v", err)
		}

		walletName := req.CompanyName + "_" + req.AdminName
		// Check if the wallet contains the identity for the provided username
		if !wallet.Exists(walletName) {
			err = addUserToWallet(wallet, walletName)
			if err != nil {
				log.Println("Failed to populate wallet contents: %v", err)
			}
		}

	}

	return &fabric_ebl.CreateCompanyResp{
		Id: userId,
	}, nil
}

func reqHttp2Rpc(req *fabric_ebl.CreateEblReq) *fabric_ipfs.CreateEblDocxReq {
	return &fabric_ipfs.CreateEblDocxReq{
		Ebl: &fabric_ipfs.CreateEblDocx{
			EblNo:                  req.Ebl.EblNo,
			OriginCompanyID:        req.Ebl.OriginCompanyID,
			OriginCompanyName:      req.Ebl.OriginCompanyName,
			ShipperCompanyID:       req.Ebl.ShipperCompanyID,
			ShipperCompanyName:     req.Ebl.ShipperCompanyName,
			ConsigneeCompanyID:     req.Ebl.ConsigneeCompanyID,
			ConsigneeCompanyName:   req.Ebl.ConsigneeCompanyName,
			NotifyPartyCompanyID:   req.Ebl.NotifyPartyCompanyID,
			NotifyPartyCompanyName: req.Ebl.NotifyPartyCompanyName,
			PlaceOfReceipt:         req.Ebl.PlaceOfReceipt,
			OceanVessel:            req.Ebl.OceanVessel,
			PortOfLoading:          req.Ebl.PortOfLoading,
			PortOfDescharge:        req.Ebl.PortOfDescharge,
			PlaceOfDestination:     req.Ebl.PlaceOfDestination,
			PlaceOfDelivery:        req.Ebl.PlaceOfDelivery,
			ShippingMarkes:         req.Ebl.ShippingMarkes,
			QuantityOfPackages:     req.Ebl.QuantityOfPackages,
			KindOfPackagesGW:       req.Ebl.KindOfPackagesGW,
			KindOfPackagesM:        req.Ebl.KindOfPackagesM,
			DescriptionOfGoods:     req.Ebl.DescriptionOfGoods,
			GrossWeight:            req.Ebl.GrossWeight,
			Measurement:            req.Ebl.Measurement,
			FreightAndCharges:      req.Ebl.FreightAndCharges,
			PlaceOfIssue:           req.Ebl.PlaceOfIssue,
			DateOfIssue:            req.Ebl.DateOfIssue,
			ShippedOnBoard:         req.Ebl.ShippedOnBoard,
			DateOfIssueDeadline:    req.Ebl.DateOfIssueDeadline,
			NumOfEBL:               req.Ebl.NumOfEBL,
			CompanyName:            req.Ebl.CompanyName,
			DeliveryAgent:          req.Ebl.DeliveryAgent,
			Status:                 req.Ebl.Status,
			File:                   req.Ebl.File,
			ContractFiles:          req.Ebl.ContractFiles,
			InvoiceFiles:           req.Ebl.InvoiceFiles,
			TransferCompanyID:      req.Ebl.TransferCompanyID,
			TransferCompanyName:    req.Ebl.TransferCompanyName,
			CompanyID:              req.Ebl.CompanyID,
		},
	}
}

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

func generateSelectorString(req *fabric_ebl.QueryEblListReq) string {
	types := 0
	selector := "{\"selector\":{"
	if req.EblFilter.EblNo != "" {
		selector += "\"eblNo\":\"" + req.EblFilter.EblNo + "\""
		types++
	}
	if req.EblFilter.OriginCompanyID != "" {
		if types > 0 {
			selector += ","
		}
		selector += "\"originCompanyID\":\"" + req.EblFilter.OriginCompanyID + "\""
		types++
	}
	if req.EblFilter.ShipperCompanyID != "" {
		if types > 0 {
			selector += ","
		}
		selector += "\"shipperCompanyID\":\"" + req.EblFilter.ShipperCompanyID + "\""
		types++
	}
	if req.EblFilter.ConsigneeCompanyID != "" {
		if types > 0 {
			selector += ","
		}
		selector += "\"consigneeCompanyID\":\"" + req.EblFilter.ConsigneeCompanyID + "\""
		types++
	}
	if req.EblFilter.NotifyPartyCompanyID != "" {
		if types > 0 {
			selector += ","
		}
		selector += "\"notifyPartyCompanyID\":\"" + req.EblFilter.NotifyPartyCompanyID + "\""
		types++
	}
	if req.EblFilter.PortOfDescharge != "" {
		if types > 0 {
			selector += ","
		}
		selector += "\"portOfDescharge\":\"" + req.EblFilter.PortOfDescharge + "\""
		types++
	}
	if req.EblFilter.Status != "" {
		if types > 0 {
			selector += ","
		}
		selector += "\"status\":\"" + req.EblFilter.Status + "\""
		types++
	}
	if req.EblFilter.CompanyID != 0 {
		if types > 0 {
			selector += ","
		}
		selector += "\"companyID\":" + strconv.FormatInt(req.EblFilter.CompanyID, 10)
		types++
	}
	if req.EblFilter.TransferCompanyID != "" {
		if types > 0 {
			selector += ","
		}
		selector += "\"transferCompanyID\":\"" + req.EblFilter.TransferCompanyID + "\""
		types++
	}
	if req.EblFilter.DateOfIssue != 0 {
		if types > 0 {
			selector += ","
		}
		selector += "\"dateOfIssue\":" + strconv.FormatInt(req.EblFilter.DateOfIssue, 10)
		types++
	}
	if req.EblFilter.ShippedOnBoard != 0 {
		if types > 0 {
			selector += ","
		}
		selector += "\"shippedOnBoard\":" + strconv.FormatInt(req.EblFilter.ShippedOnBoard, 10)
		types++
	}
	if req.EblFilter.DateOfIssueDeadline != 0 {
		if types > 0 {
			selector += ","
		}
		selector += "\"dateOfIssueDeadline\":" + strconv.FormatInt(req.EblFilter.DateOfIssueDeadline, 10)
		types++
	}
	if req.EblFilter.QuantityOfPackages != 0 {
		if types > 0 {
			selector += ","
		}
		selector += "\"quantityOfPackages\":" + strconv.FormatFloat(req.EblFilter.QuantityOfPackages, 'f', -1, 64)
		types++
	}
	if req.EblFilter.GrossWeight != 0 {
		if types > 0 {
			selector += ","
		}
		selector += "\"grossWeight\":" + strconv.FormatFloat(req.EblFilter.GrossWeight, 'f', -1, 64)
		types++
	}
	if req.EblFilter.Measurement != 0 {
		if types > 0 {
			selector += ","
		}
		selector += "\"measurement\":" + strconv.FormatFloat(req.EblFilter.Measurement, 'f', -1, 64)
		types++
	}
	if req.EblFilter.NumOfEBL != 0 {
		if types > 0 {
			selector += ","
		}
		selector += "\"numOfEbl\":" + strconv.FormatInt(req.EblFilter.NumOfEBL, 10)
		types++
	}
	if req.EblFilter.PlaceOfDelivery != "" {
		if types > 0 {
			selector += ","
		}
		selector += "\"placeOfDelivery\":\"" + req.EblFilter.PlaceOfDelivery + "\""
		types++
	}
	if req.EblFilter.PlaceOfDestination != "" {
		if types > 0 {
			selector += ","
		}
		selector += "\"placeOfDestination\":\"" + req.EblFilter.PlaceOfDestination + "\""
		types++
	}
	if req.EblFilter.PlaceOfIssue != "" {
		if types > 0 {
			selector += ","
		}
		selector += "\"placeOfIssue\":\"" + req.EblFilter.PlaceOfIssue + "\""
		types++
	}
	if req.EblFilter.PlaceOfReceipt != "" {
		if types > 0 {
			selector += ","
		}
		selector += "\"placeOfReceipt\":\"" + req.EblFilter.PlaceOfReceipt + "\""
		types++
	}
	if req.EblFilter.PortOfLoading != "" {
		if types > 0 {
			selector += ","
		}
		selector += "\"portOfLoading\":\"" + req.EblFilter.PortOfLoading + "\""
		types++
	}
	if req.EblFilter.ShippingMarkes != "" {
		if types > 0 {
			selector += ","
		}
		selector += "\"shippingMarkes\":\"" + req.EblFilter.ShippingMarkes + "\""
		types++
	}
	if req.EblFilter.FreightAndCharges != "" {
		if types > 0 {
			selector += ","
		}
		selector += "\"freightAndCharges\":\"" + req.EblFilter.FreightAndCharges + "\""
		types++

	}
	if req.EblFilter.DescriptionOfGoods != "" {
		if types > 0 {
			selector += ","
		}
		selector += "\"descriptionOfGoods\":\"" + req.EblFilter.DescriptionOfGoods + "\""
		types++
	}
	if req.EblFilter.DeliveryAgent != "" {
		if types > 0 {
			selector += ","
		}
		selector += "\"deliveryAgent\":\"" + req.EblFilter.DeliveryAgent + "\""
		types++
	}
	selector += "},\"use_index\":[\"_design/indexEblDoc\",\"indexEbl\"]}"
	return selector
}

func GetEblByRangeWithPaginationResp2DO(result []byte) (*fabric_ebl.QueryAllEblListResp, error) {
	var eblListResp GetEblByRangeWithPaginationResp
	err := json.Unmarshal(result, &eblListResp)
	if err != nil {
		return nil, err
	}
	var eblList []*fabric_ebl.Ebl
	for _, ebl := range eblListResp.Records {
		eblList = append(eblList, &fabric_ebl.Ebl{
			EblNo:                  ebl.EblNo,
			OriginCompanyID:        ebl.OriginCompanyID,
			OriginCompanyName:      ebl.OriginCompanyName,
			ShipperCompanyID:       ebl.ShipperCompanyID,
			ShipperCompanyName:     ebl.ShipperCompanyName,
			ConsigneeCompanyID:     ebl.ConsigneeCompanyID,
			ConsigneeCompanyName:   ebl.ConsigneeCompanyName,
			NotifyPartyCompanyID:   ebl.NotifyPartyCompanyID,
			NotifyPartyCompanyName: ebl.NotifyPartyCompanyName,
			PlaceOfReceipt:         ebl.PlaceOfReceipt,
			OceanVessel:            ebl.OceanVessel,
			PortOfLoading:          ebl.PortOfLoading,
			PortOfDescharge:        ebl.PortOfDescharge,
			PlaceOfDestination:     ebl.PlaceOfDestination,
			PlaceOfDelivery:        ebl.PlaceOfDelivery,
			ShippingMarkes:         ebl.ShippingMarkes,
			QuantityOfPackages:     ebl.QuantityOfPackages,
			KindOfPackagesGW:       ebl.KindOfPackagesGW,
			KindOfPackagesM:        ebl.KindOfPackagesM,
			DescriptionOfGoods:     ebl.DescriptionOfGoods,
			GrossWeight:            ebl.GrossWeight,
			Measurement:            ebl.Measurement,
			FreightAndCharges:      ebl.FreightAndCharges,
			PlaceOfIssue:           ebl.PlaceOfIssue,
			DateOfIssue:            ebl.DateOfIssue,
			DeliveryAgent:          ebl.DeliveryAgent,
			ShippedOnBoard:         ebl.ShippedOnBoard,
			NumOfEBL:               ebl.NumOfEBL,
			DateOfIssueDeadline:    ebl.DateOfIssueDeadline,
			Status:                 ebl.Status,
			File:                   ebl.File,
			ContractFiles:          strings.Split(ebl.ContractFiles, ";"),
			InvoiceFiles:           strings.Split(ebl.InvoiceFiles, ";"),
			TransferCompanyID:      ebl.TransferCompanyID,
			TransferCompanyName:    ebl.TransferCompanyName,
			CompanyID:              ebl.CompanyID,
			CompanyName:            ebl.CompanyName,
		})
	}
	return &fabric_ebl.QueryAllEblListResp{
		EblList:             eblList,
		Bookmark:            eblListResp.Bookmark,
		FetchedRecordsCount: eblListResp.FetchedRecordsCount,
	}, nil
}

type Ebl struct {
	EblNo                  string  `thrift:"eblNo,1,required" frugal:"1,required,string" json:"eblNo"`
	OriginCompanyID        string  `thrift:"originCompanyID,2,required" frugal:"2,required,string" json:"originCompanyID"`
	OriginCompanyName      string  `thrift:"originCompanyName,3,required" frugal:"3,required,string" json:"originCompanyName"`
	ShipperCompanyID       string  `thrift:"shipperCompanyID,4,required" frugal:"4,required,string" json:"shipperCompanyID"`
	ShipperCompanyName     string  `thrift:"shipperCompanyName,5,required" frugal:"5,required,string" json:"shipperCompanyName"`
	ConsigneeCompanyID     string  `thrift:"consigneeCompanyID,6,required" frugal:"6,required,string" json:"consigneeCompanyID"`
	ConsigneeCompanyName   string  `thrift:"consigneeCompanyName,7,required" frugal:"7,required,string" json:"consigneeCompanyName"`
	NotifyPartyCompanyID   string  `thrift:"notifyPartyCompanyID,8,required" frugal:"8,required,string" json:"notifyPartyCompanyID"`
	NotifyPartyCompanyName string  `thrift:"notifyPartyCompanyName,9,required" frugal:"9,required,string" json:"notifyPartyCompanyName"`
	PlaceOfReceipt         string  `thrift:"placeOfReceipt,10,required" frugal:"10,required,string" json:"placeOfReceipt"`
	OceanVessel            string  `thrift:"oceanVessel,11,required" frugal:"11,required,string" json:"oceanVessel"`
	PortOfLoading          string  `thrift:"portOfLoading,12,required" frugal:"12,required,string" json:"portOfLoading"`
	PortOfDescharge        string  `thrift:"portOfDescharge,13,required" frugal:"13,required,string" json:"portOfDescharge"`
	PlaceOfDestination     string  `thrift:"placeOfDestination,14,required" frugal:"14,required,string" json:"placeOfDestination"`
	PlaceOfDelivery        string  `thrift:"placeOfDelivery,15,required" frugal:"15,required,string" json:"placeOfDelivery"`
	ShippingMarkes         string  `thrift:"shippingMarkes,16,required" frugal:"16,required,string" json:"shippingMarkes"`
	QuantityOfPackages     float64 `thrift:"quantityOfPackages,17,required" frugal:"17,required,double" json:"quantityOfPackages"`
	KindOfPackagesGW       string  `thrift:"kindOfPackagesGW,18,required" frugal:"18,required,string" json:"kindOfPackagesGW"`
	KindOfPackagesM        string  `thrift:"kindOfPackagesM,19,required" frugal:"19,required,string" json:"kindOfPackagesM"`
	DescriptionOfGoods     string  `thrift:"descriptionOfGoods,20,required" frugal:"20,required,string" json:"descriptionOfGoods"`
	GrossWeight            float64 `thrift:"grossWeight,21,required" frugal:"21,required,double" json:"grossWeight"`
	Measurement            float64 `thrift:"measurement,22,required" frugal:"22,required,double" json:"measurement"`
	FreightAndCharges      string  `thrift:"freightAndCharges,23,required" frugal:"23,required,string" json:"freightAndCharges"`
	PlaceOfIssue           string  `thrift:"placeOfIssue,24,required" frugal:"24,required,string" json:"placeOfIssue"`
	DateOfIssue            int64   `thrift:"dateOfIssue,25,required" frugal:"25,required,i64" json:"dateOfIssue"`
	DeliveryAgent          string  `thrift:"deliveryAgent,26,required" frugal:"26,required,string" json:"deliveryAgent"`
	ShippedOnBoard         int64   `thrift:"shippedOnBoard,27,required" frugal:"27,required,i64" json:"shippedOnBoard"`
	NumOfEBL               int64   `thrift:"numOfEBL,28,required" frugal:"28,required,i64" json:"numOfEBL"`
	DateOfIssueDeadline    int64   `thrift:"dateOfIssueDeadline,29,required" frugal:"29,required,i64" json:"dateOfIssueDeadline"`
	Status                 string  `thrift:"status,30,required" frugal:"30,required,string" json:"status"`
	File                   string  `thrift:"file,31,required" frugal:"31,required,string" json:"file"`
	ContractFiles          string  `thrift:"contractFiles,32,required" frugal:"32,required,string" json:"contractFiles"`
	InvoiceFiles           string  `thrift:"invoiceFiles,33,required" frugal:"33,required,string" json:"invoiceFiles"`
	TransferCompanyID      string  `thrift:"transferCompanyID,34,required" frugal:"34,required,string" json:"transferCompanyID"`
	TransferCompanyName    string  `thrift:"transferCompanyName,35,required" frugal:"35,required,string" json:"transferCompanyName"`
	CompanyID              int64   `thrift:"companyID,36,required" frugal:"36,required,i64" json:"companyID"`
	CompanyName            string  `thrift:"companyName,37,required" frugal:"37,required,string" json:"companyName"`
}

type GetEblByRangeWithPaginationResp struct {
	Records             []*Ebl `thrift:"records,1,required" frugal:"1,required" json:"records"`
	FetchedRecordsCount int64  `thrift:"fetchedRecordsCount,2,required" frugal:"2,required,string" json:"fetchedRecordsCount"`
	Bookmark            string `thrift:"bookmark,3,required" frugal:"3,required,string" json:"bookmark"`
}
