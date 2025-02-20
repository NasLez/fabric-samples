package service

import (
	"context"
	"fabric_ebl/biz_error"
	"fabric_ebl/common/id_gen"
	"fabric_ebl/domain/converter"
	"fabric_ebl/repo"
	"fabric_ebl/sal/jwt"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/wxl-server/idl_gen/kitex_gen/fabric_ebl"
	"io/ioutil"
	"log"
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

func (u FabricEblServiceImpl) CreateEbl(ctx context.Context, req *fabric_ebl.CreateEblReq) (*fabric_ebl.CreateEblResp, error) {
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
	originCompanyID, err := strconv.ParseInt(req.Ebl.OriginCompanyID, 10, 64)
	if company.ID != originCompanyID {
		logger.CtxErrorf(ctx, "CreateEbl failed, company id is wrong")
		return nil, err
	}
	if company.Name != req.Ebl.OriginCompanyName {
		logger.CtxErrorf(ctx, "CreateEbl failed, company name is wrong")
		return nil, err
	}
	walletName := company.Name + "_" + user.Name
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}

	if !wallet.Exists(walletName) {
		err = addUserToWallet(wallet, walletName)
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
		gateway.WithIdentity(wallet, walletName),
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

	log.Println("--> Submit Transaction: CreateEbl, creates new EBL with provided details")
	ID, err := id_gen.NextID()
	req.Ebl.EblNo = strconv.FormatInt(ID, 10) + "-" + strconv.FormatInt(req.Ebl.CompanyID, 10)

	{
		result, err := contract.SubmitTransaction(
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
	}
	result, err := contract.SubmitTransaction(
		"CreateEbl",                                                  // chaincode method
		req.Ebl.EblNo,                                                // eblNo
		req.Ebl.OriginCompanyID,                                      // originCompanyID
		req.Ebl.OriginCompanyName,                                    // originCompanyName
		req.Ebl.ShipperCompanyID,                                     // shipperCompanyID
		req.Ebl.ShipperCompanyName,                                   // shipperCompanyName
		req.Ebl.ConsigneeCompanyID,                                   // consigneeCompanyID
		req.Ebl.ConsigneeCompanyName,                                 // consigneeCompanyName
		req.Ebl.NotifyPartyCompanyID,                                 // notifyPartyCompanyID
		req.Ebl.NotifyPartyCompanyName,                               // notifyPartyCompanyName
		req.Ebl.PlaceOfReceipt,                                       // placeOfReceipt
		req.Ebl.OceanVessel,                                          // oceanVessel
		req.Ebl.PortOfLoading,                                        // portOfLoading
		req.Ebl.PortOfDescharge,                                      // portOfDescharge
		req.Ebl.PlaceOfDestination,                                   // placeOfDestination
		req.Ebl.PlaceOfDelivery,                                      // placeOfDelivery
		req.Ebl.ShippingMarkes,                                       // shippingMarkes
		strings.Join(req.Ebl.ContractFiles, ";"),                     // contractFiles (can be a file or file path)
		strings.Join(req.Ebl.InvoiceFiles, ";"),                      // invoiceFiles (can be a file or file path)
		req.Ebl.TransferCompanyID,                                    // transferCompanyID
		req.Ebl.TransferCompanyName,                                  // transferCompanyName
		req.Ebl.KindOfPackagesGW,                                     // kindOfPackagesGW
		req.Ebl.KindOfPackagesM,                                      // kindOfPackagesM
		req.Ebl.DescriptionOfGoods,                                   // descriptionOfGoods
		req.Ebl.DeliveryAgent,                                        // deliveryAgent
		req.Ebl.CompanyName,                                          // companyName
		req.Ebl.FreightAndCharges,                                    // freightAndCharges
		req.Ebl.Status,                                               // status
		req.Ebl.File,                                                 // file
		req.Ebl.PlaceOfIssue,                                         // placeOfIssue
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
		log.Fatalf("Failed to Submit transaction: %v", err)
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
			log.Fatalf("Failed to create wallet: %v", err)
		}

		walletName := req.CompanyName + "_" + req.AdminName
		// Check if the wallet contains the identity for the provided username
		if !wallet.Exists(walletName) {
			err = addUserToWallet(wallet, walletName)
			if err != nil {
				log.Fatalf("Failed to populate wallet contents: %v", err)
			}
		}

	}

	return &fabric_ebl.CreateCompanyResp{
		Id: userId,
	}, nil
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
