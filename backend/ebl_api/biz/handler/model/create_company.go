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

// CreateCompany .
// @router /user/CreateCompany [GET]
func CreateCompany(ctx context.Context, c *app.RequestContext) {
	NewCreateCompanyHandler(ctx, c).handle()
}

type CreateCompanyHandler struct {
	ctx      context.Context
	hertzCtx *app.RequestContext
	respData *model.CreateCompanyData
}

func NewCreateCompanyHandler(ctx context.Context, hertzCtx *app.RequestContext) *CreateCompanyHandler {
	return &CreateCompanyHandler{
		ctx:      ctx,
		hertzCtx: hertzCtx,
	}
}

func (h *CreateCompanyHandler) handle() {
	// 业务逻辑
	ctx := h.ctx
	var req model.CreateCompanyReq
	err := h.hertzCtx.BindAndValidate(&req)
	if err != nil {
		h.hertzCtx.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp, err := fabric_ebl_rpc.CreateCompany(ctx, h.reqHttp2Rpc(&req))
	if err != nil {
		h.ReturnResp(Status.CreateCompanyError, err)
		return
	}
	h.respData = h.respRpc2Http(resp)
	h.ReturnResp(Status.Success, nil)

}
func (h *CreateCompanyHandler) respRpc2Http(resp *fabric_ebl.CreateCompanyResp) *model.CreateCompanyData {
	return &model.CreateCompanyData{
		Id: &resp.Id,
	}
}
func (h *CreateCompanyHandler) reqHttp2Rpc(req *model.CreateCompanyReq) *fabric_ebl.CreateCompanyReq {
	return &fabric_ebl.CreateCompanyReq{
		CompanyCode:   *req.CompanyCode,
		CompanyName:   *req.CompanyName,
		CompanyType:   fabric_ebl.CompanyType(*req.CompanyType),
		AdminEmail:    *req.AdminEmail,
		AdminPassword: *req.AdminPassword,
		AdminName:     *req.AdminName,
	}
}

func (h *CreateCompanyHandler) ReturnResp(status *Status.Status, err error) {
	if err != nil {
		logger.CtxErrorf(h.ctx, "CreateCompany failed, err = %v", err)
	}
	resp := new(model.CreateCompanyResp)
	resp.Code = status.Code()
	resp.Message = status.Message()
	if status.Code() == Status.Success.Code() && err == nil {
		resp.Data = h.respData
	}
	h.hertzCtx.JSON(consts.StatusOK, &resp)
}
