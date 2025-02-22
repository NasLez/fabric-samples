package model

import (
	"context"
	"ebl_api/biz/common/Status"
	"ebl_api/biz/model"
	"ebl_api/biz/sal/rpc/fabric_ebl_rpc"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/wxl-server/idl_gen/kitex_gen/fabric_ebl"
)

// UploadSeal .
// @router /user/UploadSeal [GET]
func UploadSeal(ctx context.Context, c *app.RequestContext) {
	NewUploadSealHandler(ctx, c).handle()
}

type UploadSealHandler struct {
	ctx      context.Context
	hertzCtx *app.RequestContext
	respData *model.UploadSealData
}

func NewUploadSealHandler(ctx context.Context, hertzCtx *app.RequestContext) *UploadSealHandler {
	return &UploadSealHandler{
		ctx:      ctx,
		hertzCtx: hertzCtx,
	}
}

func (h *UploadSealHandler) handle() {
	// 业务逻辑
	ctx := h.ctx
	var req model.UploadSealReq
	err := h.hertzCtx.BindAndValidate(&req)
	if err != nil {
		h.hertzCtx.String(consts.StatusBadRequest, err.Error())
		return
	}
	file, _ := h.hertzCtx.FormFile("file")
	fmt.Println(file.Filename)

	// Upload the file to specific dst
	err = h.hertzCtx.SaveUploadedFile(file, fmt.Sprintf("./file/upload/%s", file.Filename))
	if err != nil {
		h.hertzCtx.String(consts.StatusInternalServerError, err.Error())
		return
	}

	h.hertzCtx.String(consts.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	header := h.hertzCtx.Request.Header
	token := header.Get("token")
	f, err := file.Open()
	if err != nil {
		logger.CtxErrorf(ctx, "file.Open failed, err = %v", err)
		h.ReturnResp(Status.UploadSealError, err)
		return
	}
	p := make([]byte, file.Size)
	_, err = f.Read(p)
	if err != nil {
		logger.CtxErrorf(ctx, "f.Read failed, err = %v", err)
		h.ReturnResp(Status.UploadSealError, err)
		return
	}
	resp, err := fabric_ebl_rpc.UploadSeal(ctx, h.reqHttp2Rpc(token, p))
	if err != nil {
		h.ReturnResp(Status.UploadSealError, err)
		return
	}
	h.respData = h.respRpc2Http(resp)
	h.ReturnResp(Status.Success, nil)

}

func (h *UploadSealHandler) respRpc2Http(resp *fabric_ebl.UploadSealResp) *model.UploadSealData {
	return &model.UploadSealData{
		Id: &resp.Id,
	}
}
func byteSliceToInt8Slice(b []byte) []int8 {
	result := make([]int8, len(b))
	for i, v := range b {
		result[i] = int8(v)
	}
	return result
}

func (h *UploadSealHandler) reqHttp2Rpc(token string, p []byte) *fabric_ebl.UploadSealReq {
	// 将 []byte 转换为 []int8
	pInt8 := byteSliceToInt8Slice(p)
	return &fabric_ebl.UploadSealReq{
		Token: token,
		Seal:  pInt8,
	}
}

func (h *UploadSealHandler) ReturnResp(status *Status.Status, err error) {
	if err != nil {
		logger.CtxErrorf(h.ctx, "UploadSeal failed, err = %v", err)
	}
	resp := new(model.UploadSealResp)
	resp.Code = status.Code()
	resp.Message = status.Message()
	if status.Code() == Status.Success.Code() && err == nil {
		resp.Data = h.respData
	}
	h.hertzCtx.JSON(consts.StatusOK, &resp)
}
