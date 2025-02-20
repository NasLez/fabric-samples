package model

import (
	"context"
	"ebl_api/biz/common/Status"
	"ebl_api/biz/model"
	"ebl_api/biz/sal/rpc/fabric_ebl_rpc"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/wxl-server/idl_gen/kitex_gen/fabric_ebl"
)

// CreateEbl .
// @router /user/CreateEbl [GET]
func CreateEbl(ctx context.Context, c *app.RequestContext) {
	NewCreateEblHandler(ctx, c).handle()
}

type CreateEblHandler struct {
	ctx      context.Context
	hertzCtx *app.RequestContext
	respData *model.CreateEblData
}

func NewCreateEblHandler(ctx context.Context, hertzCtx *app.RequestContext) *CreateEblHandler {
	return &CreateEblHandler{
		ctx:      ctx,
		hertzCtx: hertzCtx,
	}
}

func (h *CreateEblHandler) handle() {
	// 业务逻辑
	ctx := h.ctx
	var req model.CreateEblReq
	err := h.hertzCtx.BindAndValidate(&req)
	if err != nil {
		h.hertzCtx.String(consts.StatusBadRequest, err.Error())
		return
	}
	header := h.hertzCtx.Request.Header
	token := header.Get("token")
	resp, err := fabric_ebl_rpc.CreateEbl(ctx, h.reqHttp2Rpc(&req, token))
	if err != nil {
		h.ReturnResp(Status.CreateEblError, err)
		return
	}
	h.respData = h.respRpc2Http(resp)
	h.ReturnResp(Status.Success, nil)

}
func (h *CreateEblHandler) respRpc2Http(resp *fabric_ebl.CreateEblResp) *model.CreateEblData {
	return &model.CreateEblData{
		Id: &resp.Id,
	}
}
func (h *CreateEblHandler) reqHttp2Rpc(req *model.CreateEblReq, token string) *fabric_ebl.CreateEblReq {
	return &fabric_ebl.CreateEblReq{
		Ebl: &fabric_ebl.Ebl{
			EblNo:                  *req.Ebl.EblNo,
			OriginCompanyID:        *req.Ebl.OriginCompanyID,
			OriginCompanyName:      *req.Ebl.OriginCompanyName,
			ShipperCompanyID:       *req.Ebl.ShipperCompanyID,
			ShipperCompanyName:     *req.Ebl.ShipperCompanyName,
			ConsigneeCompanyID:     *req.Ebl.ConsigneeCompanyID,
			ConsigneeCompanyName:   *req.Ebl.ConsigneeCompanyName,
			NotifyPartyCompanyID:   *req.Ebl.NotifyPartyCompanyID,
			NotifyPartyCompanyName: *req.Ebl.NotifyPartyCompanyName,
			PlaceOfReceipt:         *req.Ebl.PlaceOfReceipt,
			OceanVessel:            *req.Ebl.OceanVessel,
			PortOfLoading:          *req.Ebl.PortOfLoading,
			PortOfDescharge:        *req.Ebl.PortOfDescharge,
			PlaceOfDestination:     *req.Ebl.PlaceOfDestination,
			PlaceOfDelivery:        *req.Ebl.PlaceOfDelivery,
			ShippingMarkes:         *req.Ebl.ShippingMarkes,
			QuantityOfPackages:     *req.Ebl.QuantityOfPackages,
			KindOfPackagesGW:       *req.Ebl.KindOfPackagesGW,
			KindOfPackagesM:        *req.Ebl.KindOfPackagesM,
			DescriptionOfGoods:     *req.Ebl.DescriptionOfGoods,
			GrossWeight:            *req.Ebl.GrossWeight,
			Measurement:            *req.Ebl.Measurement,
			FreightAndCharges:      *req.Ebl.FreightAndCharges,
			PlaceOfIssue:           *req.Ebl.PlaceOfIssue,
			DateOfIssue:            *req.Ebl.DateOfIssue,
			DeliveryAgent:          *req.Ebl.DeliveryAgent,
			ShippedOnBoard:         *req.Ebl.ShippedOnBoard,
			NumOfEBL:               *req.Ebl.NumOfEBL,
			DateOfIssueDeadline:    *req.Ebl.DateOfIssueDeadline,
			Status:                 *req.Ebl.Status,
			File:                   *req.Ebl.File,
			ContractFiles:          req.Ebl.ContractFiles,
			InvoiceFiles:           req.Ebl.InvoiceFiles,
			TransferCompanyID:      *req.Ebl.TransferCompanyID,
			TransferCompanyName:    *req.Ebl.TransferCompanyName,
			CompanyID:              *req.Ebl.CompanyID,
			CompanyName:            *req.Ebl.CompanyName,
		},
		Token: token,
	}
}

func (h *CreateEblHandler) ReturnResp(status *Status.Status, err error) {
	if err != nil {
		logger.CtxErrorf(h.ctx, "CreateEbl failed, err = %v", err)
	}
	resp := new(model.CreateEblResp)
	resp.Code = status.Code()
	resp.Message = status.Message()
	if status.Code() == Status.Success.Code() && err == nil {
		resp.Data = h.respData
	}
	h.hertzCtx.JSON(consts.StatusOK, &resp)
}
