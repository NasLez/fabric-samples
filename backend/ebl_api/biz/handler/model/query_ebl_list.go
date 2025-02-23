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

// QueryEblList .
// @router /user/QueryEblList [GET]
func QueryEblList(ctx context.Context, c *app.RequestContext) {
	NewQueryEblListHandler(ctx, c).handle()
}

type QueryEblListHandler struct {
	ctx      context.Context
	hertzCtx *app.RequestContext
	respData *model.QueryEblListData
}

func NewQueryEblListHandler(ctx context.Context, hertzCtx *app.RequestContext) *QueryEblListHandler {
	return &QueryEblListHandler{
		ctx:      ctx,
		hertzCtx: hertzCtx,
	}
}

func (h *QueryEblListHandler) handle() {
	// 业务逻辑
	ctx := h.ctx
	var req model.QueryEblListReq
	err := h.hertzCtx.BindAndValidate(&req)
	if err != nil {
		h.hertzCtx.String(consts.StatusBadRequest, err.Error())
		return
	}
	header := h.hertzCtx.Request.Header
	token := header.Get("token")
	resp, err := fabric_ebl_rpc.QueryEblList(ctx, h.reqHttp2Rpc(&req, token))
	if err != nil {
		h.ReturnResp(Status.QueryEblListError, err)
		return
	}
	h.respData = h.respRpc2Http(resp)
	h.ReturnResp(Status.Success, nil)

}
func (h *QueryEblListHandler) respRpc2Http(resp *fabric_ebl.QueryEblListResp) *model.QueryEblListData {
	eblList := make([]*model.Ebl, 0)
	for _, v := range resp.EblList {
		eblList = append(eblList, &model.Ebl{
			EblNo:                  &v.EblNo,
			OriginCompanyID:        &v.OriginCompanyID,
			OriginCompanyName:      &v.OriginCompanyName,
			ShipperCompanyID:       &v.ShipperCompanyID,
			ShipperCompanyName:     &v.ShipperCompanyName,
			ConsigneeCompanyID:     &v.ConsigneeCompanyID,
			ConsigneeCompanyName:   &v.ConsigneeCompanyName,
			NotifyPartyCompanyID:   &v.NotifyPartyCompanyID,
			NotifyPartyCompanyName: &v.NotifyPartyCompanyName,
			PlaceOfReceipt:         &v.PlaceOfReceipt,
			OceanVessel:            &v.OceanVessel,
			PortOfLoading:          &v.PortOfLoading,
			PortOfDescharge:        &v.PortOfDescharge,
			PlaceOfDestination:     &v.PlaceOfDestination,
			PlaceOfDelivery:        &v.PlaceOfDelivery,
			ShippingMarkes:         &v.ShippingMarkes,
			QuantityOfPackages:     &v.QuantityOfPackages,
			KindOfPackagesGW:       &v.KindOfPackagesGW,
			KindOfPackagesM:        &v.KindOfPackagesM,
			DescriptionOfGoods:     &v.DescriptionOfGoods,
			GrossWeight:            &v.GrossWeight,
			Measurement:            &v.Measurement,
			FreightAndCharges:      &v.FreightAndCharges,
			PlaceOfIssue:           &v.PlaceOfIssue,
			DateOfIssue:            &v.DateOfIssue,
			DeliveryAgent:          &v.DeliveryAgent,
			ShippedOnBoard:         &v.ShippedOnBoard,
			NumOfEBL:               &v.NumOfEBL,
			DateOfIssueDeadline:    &v.DateOfIssueDeadline,
			Status:                 &v.Status,
			File:                   &v.File,
			ContractFiles:          v.ContractFiles,
			InvoiceFiles:           v.InvoiceFiles,
			TransferCompanyID:      &v.TransferCompanyID,
			TransferCompanyName:    &v.TransferCompanyName,
			CompanyID:              &v.CompanyID,
			CompanyName:            &v.CompanyName,
		})
	}
	return &model.QueryEblListData{
		Bookmark:            &resp.Bookmark,
		EblList:             eblList,
		FetchedRecordsCount: &resp.FetchedRecordsCount,
	}
}
func (h *QueryEblListHandler) reqHttp2Rpc(req *model.QueryEblListReq, token string) *fabric_ebl.QueryEblListReq {

	return &fabric_ebl.QueryEblListReq{
		Token:    token,
		PageSize: req.PageSize,
		Bookmark: req.Bookmark,
		EblFilter: &fabric_ebl.EblFilter{
			EblNo:                *req.Ebl.EblNo,
			OriginCompanyID:      *req.Ebl.OriginCompanyID,
			ShipperCompanyID:     *req.Ebl.ShipperCompanyID,
			ConsigneeCompanyID:   *req.Ebl.ConsigneeCompanyID,
			NotifyPartyCompanyID: *req.Ebl.NotifyPartyCompanyID,
			PlaceOfReceipt:       *req.Ebl.PlaceOfReceipt,
			OceanVessel:          *req.Ebl.OceanVessel,
			PortOfLoading:        *req.Ebl.PortOfLoading,
			PortOfDescharge:      *req.Ebl.PortOfDescharge,
			PlaceOfDestination:   *req.Ebl.PlaceOfDestination,
			PlaceOfDelivery:      *req.Ebl.PlaceOfDelivery,
			ShippingMarkes:       *req.Ebl.ShippingMarkes,
			QuantityOfPackages:   *req.Ebl.QuantityOfPackages,
			KindOfPackagesGW:     *req.Ebl.KindOfPackagesGW,
			KindOfPackagesM:      *req.Ebl.KindOfPackagesM,
			DescriptionOfGoods:   *req.Ebl.DescriptionOfGoods,
			GrossWeight:          *req.Ebl.GrossWeight,
			Measurement:          *req.Ebl.Measurement,
			FreightAndCharges:    *req.Ebl.FreightAndCharges,
			PlaceOfIssue:         *req.Ebl.PlaceOfIssue,
			DateOfIssue:          *req.Ebl.DateOfIssue,
			DeliveryAgent:        *req.Ebl.DeliveryAgent,
			ShippedOnBoard:       *req.Ebl.ShippedOnBoard,
			NumOfEBL:             *req.Ebl.NumOfEBL,
			DateOfIssueDeadline:  *req.Ebl.DateOfIssueDeadline,
			Status:               *req.Ebl.Status,
			TransferCompanyID:    *req.Ebl.TransferCompanyID,
			CompanyID:            *req.Ebl.CompanyID,
		},
	}
}

func (h *QueryEblListHandler) ReturnResp(status *Status.Status, err error) {
	if err != nil {
		logger.CtxErrorf(h.ctx, "QueryEblList failed, err = %v", err)
	}
	resp := new(model.QueryEblListResp)
	resp.Code = status.Code()
	resp.Message = status.Message()
	if status.Code() == Status.Success.Code() && err == nil {
		resp.Data = h.respData
	}
	h.hertzCtx.JSON(consts.StatusOK, &resp)
}
