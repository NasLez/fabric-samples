package model

import (
	"context"
	"ebl_api/biz/common/Status"
	"ebl_api/biz/model"
	"ebl_api/biz/sal/rpc/fabric_ebl_rpc"
	"github.com/wxl-server/idl_gen/kitex_gen/fabric_ebl"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// Login .
// @router /user/login [GET]
func Login(ctx context.Context, c *app.RequestContext) {
	NewLoginHandler(ctx, c).handle()
}

type LoginHandler struct {
	ctx      context.Context
	hertzCtx *app.RequestContext
	respData *model.LoginData
}

func NewLoginHandler(ctx context.Context, hertzCtx *app.RequestContext) *LoginHandler {
	return &LoginHandler{
		ctx:      ctx,
		hertzCtx: hertzCtx,
	}
}

func (h *LoginHandler) handle() {
	// 业务逻辑
	ctx := h.ctx
	var req model.LoginReq
	err := h.hertzCtx.BindAndValidate(&req)
	if err != nil {
		h.hertzCtx.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp, err := fabric_ebl_rpc.Login(ctx, h.reqHttp2Rpc(&req))
	if err != nil {
		h.ReturnResp(Status.LoginError, err)
		return
	}
	h.respData = h.respRpc2Http(resp)
	h.ReturnResp(Status.Success, nil)
}
func (h *LoginHandler) respRpc2Http(resp *fabric_ebl.LoginResp) *model.LoginData {
	return &model.LoginData{
		Token: &resp.Token,
	}
}
func (h *LoginHandler) reqHttp2Rpc(req *model.LoginReq) *fabric_ebl.LoginReq {
	return &fabric_ebl.LoginReq{
		Email:    *req.Email,
		Password: *req.Password,
	}
}

func (h *LoginHandler) ReturnResp(status *Status.Status, err error) {
	if err != nil {
		logger.CtxErrorf(h.ctx, "Login failed, err = %v", err)
	}
	resp := new(model.LoginResp)
	resp.Code = status.Code()
	resp.Message = status.Message()
	if status.Code() == Status.Success.Code() && err == nil {
		resp.Data = h.respData
	}
	h.hertzCtx.JSON(consts.StatusOK, &resp)
}
