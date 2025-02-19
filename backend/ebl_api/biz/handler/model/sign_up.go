package model

import (
	"context"
	"ebl_api/biz/model"
	"ebl_api/biz/sal/rpc/common_user_rpc"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/wxl-server/common/Status"
	"github.com/wxl-server/idl_gen/kitex_gen/common_user"
)

// SignUp .
// @router /user/signup [GET]
func SignUp(ctx context.Context, c *app.RequestContext) {
	NewSignUpHandler(ctx, c).handle()
}

type SignUpHandler struct {
	ctx      context.Context
	hertzCtx *app.RequestContext
	respData *model.SignUpData
}

func NewSignUpHandler(ctx context.Context, hertzCtx *app.RequestContext) *SignUpHandler {
	return &SignUpHandler{
		ctx:      ctx,
		hertzCtx: hertzCtx,
	}
}

func (h *SignUpHandler) handle() {
	// 业务逻辑
	ctx := h.ctx
	var req model.SignUpReq
	err := h.hertzCtx.BindAndValidate(&req)
	if err != nil {
		h.hertzCtx.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp, err := common_user_rpc.SignUp(ctx, h.reqHttp2Rpc(&req))
	if err != nil {
		return
	}
	h.respData = h.respRpc2Http(resp)
	h.ReturnResp(Status.Success, nil)
}
func (h *SignUpHandler) respRpc2Http(resp *common_user.SignUpResp) *model.SignUpData {
	return &model.SignUpData{
		Id: &resp.Id,
	}
}
func (h *SignUpHandler) reqHttp2Rpc(req *model.SignUpReq) *common_user.SignUpReq {
	return &common_user.SignUpReq{
		Email:    *req.Email,
		Password: *req.Password,
	}
}

func (h *SignUpHandler) ReturnResp(status *Status.Status, err error) {
	if err != nil {
		logger.CtxErrorf(h.ctx, "SignUp failed, err = %v", err)
	}
	resp := new(model.SignUpResp)
	resp.Code = status.Code()
	resp.Message = status.Message()
	if status.Code() == Status.Success.Code() && err == nil {
		resp.Data = h.respData
	}
	h.hertzCtx.JSON(consts.StatusOK, &resp)
}
