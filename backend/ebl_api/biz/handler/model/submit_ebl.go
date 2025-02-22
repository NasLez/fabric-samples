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

// SubmitEbl .
// @router /user/SubmitEbl [GET]
func SubmitEbl(ctx context.Context, c *app.RequestContext) {
	NewSubmitEblHandler(ctx, c).handle()
}

type SubmitEblHandler struct {
	ctx      context.Context
	hertzCtx *app.RequestContext
	respData *model.SubmitEblData
}

func NewSubmitEblHandler(ctx context.Context, hertzCtx *app.RequestContext) *SubmitEblHandler {
	return &SubmitEblHandler{
		ctx:      ctx,
		hertzCtx: hertzCtx,
	}
}

func (h *SubmitEblHandler) handle() {
	// 业务逻辑
	ctx := h.ctx
	var req model.SubmitEblReq
	err := h.hertzCtx.BindAndValidate(&req)
	if err != nil {
		h.hertzCtx.String(consts.StatusBadRequest, err.Error())
		return
	}
	header := h.hertzCtx.Request.Header
	token := header.Get("token")
	resp, err := fabric_ebl_rpc.SubmitEbl(ctx, h.reqHttp2Rpc(&req, token))
	if err != nil {
		h.ReturnResp(Status.SubmitEblError, err)
		return
	}
	h.respData = h.respRpc2Http(resp)
	h.ReturnResp(Status.Success, nil)

}
func (h *SubmitEblHandler) respRpc2Http(resp *fabric_ebl.SubmitEblResp) *model.SubmitEblData {
	return &model.SubmitEblData{
		Id: &resp.Id,
	}
}
func (h *SubmitEblHandler) reqHttp2Rpc(req *model.SubmitEblReq, token string) *fabric_ebl.SubmitEblReq {
	return &fabric_ebl.SubmitEblReq{
		EblNo: *req.EblNo,
		Token: token,
	}
}

func (h *SubmitEblHandler) ReturnResp(status *Status.Status, err error) {
	if err != nil {
		logger.CtxErrorf(h.ctx, "SubmitEbl failed, err = %v", err)
	}
	resp := new(model.SubmitEblResp)
	resp.Code = status.Code()
	resp.Message = status.Message()
	if status.Code() == Status.Success.Code() && err == nil {
		resp.Data = h.respData
	}
	h.hertzCtx.JSON(consts.StatusOK, &resp)
}
