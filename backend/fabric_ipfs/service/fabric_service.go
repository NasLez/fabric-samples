package service

import (
	"bytes"
	"context"
	"fabric_ipfs/repo"
	"fmt"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/nguyenthenguyen/docx"
	"github.com/wxl-server/idl_gen/kitex_gen/fabric_ipfs"
	"go.uber.org/dig"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type FabricIpfsService interface {
	CreateEblDocx(ctx context.Context, req *fabric_ipfs.CreateEblDocxReq) (*fabric_ipfs.CreateEblDocxResp, error)
}

type Param struct {
	dig.In
	FabricEblRepo repo.FabricIpfsRepo
}

type FabricIpfsServiceImpl struct {
	p Param
}

func (f FabricIpfsServiceImpl) CreateEblDocx(ctx context.Context, req *fabric_ipfs.CreateEblDocxReq) (*fabric_ipfs.CreateEblDocxResp, error) {
	r, err := docx.ReadDocxFile("./ebl_template.docx")
	// 或者从内存中读取
	// r, err := docx.ReadDocxFromMemory(data io.ReaderAt, size int64)

	// 或从文件系统对象读取：
	// r, err := docx.ReadDocxFromFS(file string, fs fs.FS)

	if err != nil {
		panic(err)
	}
	docx1 := r.Editable()
	// 类似于Go标准库中的strings.Replace使用方法
	docx1.Replace("eblNo", req.Ebl.EblNo, -1)
	docx1.Replace("shipper", req.Ebl.ShipperCompanyName, -1)
	docx1.Replace("consignee", req.Ebl.ConsigneeCompanyName, -1)
	docx1.Replace("notifyParty", req.Ebl.NotifyPartyCompanyName, -1)
	docx1.Replace("placeOfReceipt", req.Ebl.PlaceOfReceipt, -1)
	docx1.Replace("oceanVessel", req.Ebl.OceanVessel, -1)
	docx1.Replace("portOfLoading", req.Ebl.PortOfLoading, -1)
	docx1.Replace("portOfDescharge", req.Ebl.PortOfDescharge, -1)
	docx1.Replace("placeOfDestination", req.Ebl.PlaceOfDestination, -1)
	docx1.Replace("placeOfDelivery", req.Ebl.PlaceOfDelivery, -1)
	docx1.Replace("shippingMarkes", req.Ebl.ShippingMarkes, -1)
	docx1.Replace("quantityOfPackages", strconv.FormatFloat(req.Ebl.QuantityOfPackages, 'f', -1, 64), -1)
	docx1.Replace("kindOfPackagesGW", req.Ebl.KindOfPackagesGW, -1)
	docx1.Replace("kindOfPackagesM", req.Ebl.KindOfPackagesM, -1)
	docx1.Replace("descriptionOfGoods", req.Ebl.DescriptionOfGoods, -1)
	docx1.Replace("grossWeight", strconv.FormatFloat(req.Ebl.GrossWeight, 'f', -1, 64), -1)
	docx1.Replace("measurement", strconv.FormatFloat(req.Ebl.Measurement, 'f', -1, 64), -1)
	docx1.Replace("freightAndCharges", req.Ebl.FreightAndCharges, -1)
	docx1.Replace("placeOfIssue", req.Ebl.PlaceOfIssue, -1)
	docx1.Replace("dateOfIssue", strconv.FormatInt(req.Ebl.DateOfIssue, 10), -1)
	docx1.Replace("deliveryAgent", req.Ebl.DeliveryAgent, -1)
	docx1.Replace("shippedOnBoard", strconv.FormatInt(req.Ebl.ShippedOnBoard, 10), -1)
	docx1.Replace("numOfEbl", strconv.FormatInt(req.Ebl.NumOfEBL, 10), -1)
	docx1.Replace("dateOfIssueDeadline", strconv.FormatInt(req.Ebl.DateOfIssueDeadline, 10), -1)
	docx1.WriteToFile("./docx/" + req.Ebl.EblNo + ".docx")
	raw := Read("./docx/" + req.Ebl.EblNo + ".docx")
	resp := fabric_ipfs.CreateEblDocxResp{}
	if raw != nil {
		hash, err := UploadIPFS(raw)
		if err != nil {
			log.Println("UploadIPFS err", err)
		}
		log.Println("hash", hash)
		resp.FileHash = hash
	} else {
		log.Println("read file fail")
	}
	//delete file
	//err = os.Remove("./docx/" + req.Ebl.EblNo + ".docx")
	//if err != nil {
	//	log.Println("delete file fail", err)
	//	return nil, err
	//}
	return &resp, nil
}

func NewFabricIpfsService(p Param) FabricIpfsService {
	return &FabricIpfsServiceImpl{
		p: p,
	}
}
func Read(filepath string) []byte {
	f, err := os.Open(filepath)
	if err != nil {
		log.Println("read file fail", err)
		return nil
	}
	defer f.Close()

	fd, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println("read to fd fail", err)
		return nil
	}

	return fd
}
func UploadIPFS(raw []byte) (string, error) {
	sh := shell.NewShell("localhost:5001")
	reader := bytes.NewReader(raw)
	// https://github.com/ipfs/go-ipfs-api/blob/master/add.go
	fileHash, err := sh.Add(reader)
	if err != nil {
		return "", err
	}
	fmt.Println(fileHash)
	return fileHash, nil
}
