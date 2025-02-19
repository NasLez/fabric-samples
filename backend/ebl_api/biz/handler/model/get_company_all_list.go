package model

import (
	"context"
	"ebl_api/biz/common/Status"
	"ebl_api/biz/model"
	"ebl_api/biz/sal/rpc/fabric_ebl_rpc"
	"ebl_api/common/gptr"
	"ebl_api/common/gslice"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/wxl-server/idl_gen/kitex_gen/fabric_ebl"
)

// GetCompanyAllList .
// @router /user/GetCompanyAllList [GET]
func GetCompanyAllList(ctx context.Context, c *app.RequestContext) {
	NewGetCompanyAllListHandler(ctx, c).handle()
}

type GetCompanyAllListHandler struct {
	ctx      context.Context
	hertzCtx *app.RequestContext
	respData *model.GetCompanyAllListData
}

func NewGetCompanyAllListHandler(ctx context.Context, hertzCtx *app.RequestContext) *GetCompanyAllListHandler {
	return &GetCompanyAllListHandler{
		ctx:      ctx,
		hertzCtx: hertzCtx,
	}
}

func (h *GetCompanyAllListHandler) handle() {
	// 业务逻辑
	ctx := h.ctx
	var req model.GetCompanyAllListReq
	err := h.hertzCtx.BindAndValidate(&req)
	if err != nil {
		h.hertzCtx.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp, err := fabric_ebl_rpc.GetCompanyAllList(ctx, h.reqHttp2Rpc())
	if err != nil {
		h.ReturnResp(Status.GetCompanyAllListError, err)
		return
	}
	h.respData = h.respRpc2Http(resp)
	h.ReturnResp(Status.Success, nil)

}
func (h *GetCompanyAllListHandler) respRpc2Http(resp *fabric_ebl.GetCompanyAllListResp) *model.GetCompanyAllListData {
	return &model.GetCompanyAllListData{
		CompanyList: gslice.Map(resp.CompanyList, func(v *fabric_ebl.Company) *model.Company {
			return &model.Company{
				Id:          gptr.Of(v.Id),
				CompanyCode: gptr.Of(v.CompanyCode),
				CompanyName: gptr.Of(v.CompanyName),
				CompanyType: gptr.Of(model.CompanyType(v.CompanyType)),
			}
		}),
	}
}
func (h *GetCompanyAllListHandler) reqHttp2Rpc() *fabric_ebl.GetCompanyAllListReq {
	return &fabric_ebl.GetCompanyAllListReq{}
}

func (h *GetCompanyAllListHandler) ReturnResp(status *Status.Status, err error) {
	if err != nil {
		logger.CtxErrorf(h.ctx, "GetCompanyAllList failed, err = %v", err)
	}
	resp := new(model.GetCompanyAllListResp)
	resp.Code = status.Code()
	resp.Message = status.Message()
	if status.Code() == Status.Success.Code() && err == nil {
		resp.Data = h.respData
	}
	h.hertzCtx.JSON(consts.StatusOK, &resp)
}
