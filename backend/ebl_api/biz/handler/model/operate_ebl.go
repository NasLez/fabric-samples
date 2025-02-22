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

// OperateEbl .
// @router /user/OperateEbl [GET]
func OperateEbl(ctx context.Context, c *app.RequestContext) {
	NewOperateEblHandler(ctx, c).handle()
}

type OperateEblHandler struct {
	ctx      context.Context
	hertzCtx *app.RequestContext
	respData *model.OperateEblData
}

func NewOperateEblHandler(ctx context.Context, hertzCtx *app.RequestContext) *OperateEblHandler {
	return &OperateEblHandler{
		ctx:      ctx,
		hertzCtx: hertzCtx,
	}
}

func (h *OperateEblHandler) handle() {
	// 业务逻辑
	ctx := h.ctx
	var req model.OperateEblReq
	err := h.hertzCtx.BindAndValidate(&req)
	if err != nil {
		h.hertzCtx.String(consts.StatusBadRequest, err.Error())
		return
	}
	header := h.hertzCtx.Request.Header
	token := header.Get("token")
	resp, err := fabric_ebl_rpc.OperateEbl(ctx, h.reqHttp2Rpc(&req, token))
	if err != nil {
		h.ReturnResp(Status.OperateEblError, err)
		return
	}
	h.respData = h.respRpc2Http(resp)
	h.ReturnResp(Status.Success, nil)

}
func (h *OperateEblHandler) respRpc2Http(resp *fabric_ebl.OperateEblResp) *model.OperateEblData {
	return &model.OperateEblData{
		Id: &resp.Id,
	}
}
func (h *OperateEblHandler) reqHttp2Rpc(req *model.OperateEblReq, token string) *fabric_ebl.OperateEblReq {
	return &fabric_ebl.OperateEblReq{
		EblNo: *req.EblNo,
		Token: token,
		Type:  fabric_ebl.OperationType(*req.Type),
	}
}

func (h *OperateEblHandler) ReturnResp(status *Status.Status, err error) {
	if err != nil {
		logger.CtxErrorf(h.ctx, "OperateEbl failed, err = %v", err)
	}
	resp := new(model.OperateEblResp)
	resp.Code = status.Code()
	resp.Message = status.Message()
	if status.Code() == Status.Success.Code() && err == nil {
		resp.Data = h.respData
	}
	h.hertzCtx.JSON(consts.StatusOK, &resp)
}
