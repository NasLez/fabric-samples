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

// QueryAllEblList .
// @router /user/QueryAllEblList [GET]
func QueryAllEblList(ctx context.Context, c *app.RequestContext) {
	NewQueryAllEblListHandler(ctx, c).handle()
}

type QueryAllEblListHandler struct {
	ctx      context.Context
	hertzCtx *app.RequestContext
	respData *model.QueryAllEblListData
}

func NewQueryAllEblListHandler(ctx context.Context, hertzCtx *app.RequestContext) *QueryAllEblListHandler {
	return &QueryAllEblListHandler{
		ctx:      ctx,
		hertzCtx: hertzCtx,
	}
}

func (h *QueryAllEblListHandler) handle() {
	// 业务逻辑
	ctx := h.ctx
	var req model.QueryAllEblListReq
	err := h.hertzCtx.BindAndValidate(&req)
	if err != nil {
		h.hertzCtx.String(consts.StatusBadRequest, err.Error())
		return
	}
	header := h.hertzCtx.Request.Header
	token := header.Get("token")
	resp, err := fabric_ebl_rpc.QueryAllEblList(ctx, h.reqHttp2Rpc(&req, token))
	if err != nil {
		h.ReturnResp(Status.QueryAllEblListError, err)
		return
	}
	h.respData = h.respRpc2Http(resp)
	h.ReturnResp(Status.Success, nil)

}
func (h *QueryAllEblListHandler) respRpc2Http(resp *fabric_ebl.QueryAllEblListResp) *model.QueryAllEblListData {
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
	return &model.QueryAllEblListData{
		Bookmark:            &resp.Bookmark,
		EblList:             eblList,
		FetchedRecordsCount: &resp.FetchedRecordsCount,
	}
}
func (h *QueryAllEblListHandler) reqHttp2Rpc(req *model.QueryAllEblListReq, token string) *fabric_ebl.QueryAllEblListReq {

	return &fabric_ebl.QueryAllEblListReq{
		Token:    token,
		PageSize: req.PageSize,
		Bookmark: req.Bookmark,
	}
}

func (h *QueryAllEblListHandler) ReturnResp(status *Status.Status, err error) {
	if err != nil {
		logger.CtxErrorf(h.ctx, "QueryAllEblList failed, err = %v", err)
	}
	resp := new(model.QueryAllEblListResp)
	resp.Code = status.Code()
	resp.Message = status.Message()
	if status.Code() == Status.Success.Code() && err == nil {
		resp.Data = h.respData
	}
	h.hertzCtx.JSON(consts.StatusOK, &resp)
}
