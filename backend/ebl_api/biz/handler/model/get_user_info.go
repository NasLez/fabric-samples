package model

import (
	"context"
	"ebl_api/biz/common/Status"
	"ebl_api/biz/model"
	"ebl_api/biz/sal/rpc/fabric_ebl_rpc"
	"ebl_api/common/gptr"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/wxl-server/idl_gen/kitex_gen/fabric_ebl"
)

// GetUserInfo .
// @router /user/GetUserInfo [GET]
func GetUserInfo(ctx context.Context, c *app.RequestContext) {
	NewGetUserInfoHandler(ctx, c).handle()
}

type GetUserInfoHandler struct {
	ctx      context.Context
	hertzCtx *app.RequestContext
	respData *model.GetUserInfoData
}

func NewGetUserInfoHandler(ctx context.Context, hertzCtx *app.RequestContext) *GetUserInfoHandler {
	return &GetUserInfoHandler{
		ctx:      ctx,
		hertzCtx: hertzCtx,
	}
}

func (h *GetUserInfoHandler) handle() {
	// 业务逻辑
	ctx := h.ctx
	var req model.GetUserInfoReq
	err := h.hertzCtx.BindAndValidate(&req)
	if err != nil {
		h.hertzCtx.String(consts.StatusBadRequest, err.Error())
		return
	}
	header := h.hertzCtx.Request.Header
	token := header.Get("token")
	resp, err := fabric_ebl_rpc.GetUserInfo(ctx, h.reqHttp2Rpc(token))
	if err != nil {
		h.ReturnResp(Status.GetUserInfoError, err)
		return
	}
	h.respData = h.respRpc2Http(resp)
	h.ReturnResp(Status.Success, nil)

}
func (h *GetUserInfoHandler) respRpc2Http(resp *fabric_ebl.GetUserInfoResp) *model.GetUserInfoData {
	return &model.GetUserInfoData{
		UserId:      &resp.UserId,
		UserName:    &resp.UserName,
		UserEmail:   &resp.UserEmail,
		UserType:    gptr.Of(model.UserType(resp.UserType)),
		CompanyCode: &resp.CompanyCode,
		CompanyName: &resp.CompanyName,
		CompanyId:   &resp.CompanyId,
		CompanyType: gptr.Of(model.CompanyType(resp.CompanyType)),
	}
}
func (h *GetUserInfoHandler) reqHttp2Rpc(token string) *fabric_ebl.GetUserInfoReq {
	return &fabric_ebl.GetUserInfoReq{
		Token: token,
	}
}

func (h *GetUserInfoHandler) ReturnResp(status *Status.Status, err error) {
	if err != nil {
		logger.CtxErrorf(h.ctx, "GetUserInfo failed, err = %v", err)
	}
	resp := new(model.GetUserInfoResp)
	resp.Code = status.Code()
	resp.Message = status.Message()
	if status.Code() == Status.Success.Code() && err == nil {
		resp.Data = h.respData
	}
	h.hertzCtx.JSON(consts.StatusOK, &resp)
}
