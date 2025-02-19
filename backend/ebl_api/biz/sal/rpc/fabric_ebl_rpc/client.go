package fabric_ebl_rpc

import (
	"context"
	"ebl_api/common/wxl_cluster"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/wxl-server/idl_gen/kitex_gen/fabric_ebl"
	"github.com/wxl-server/idl_gen/kitex_gen/fabric_ebl/fabricebl"
)

var client = wxl_cluster.NewClient(fabricebl.NewClient, "fabric_ebl")

func CreateCompany(ctx context.Context, req *fabric_ebl.CreateCompanyReq) (resp *fabric_ebl.CreateCompanyResp, err error) {
	resp, err = client.CreateCompany(ctx, req)
	if err != nil {
		logger.CtxErrorf(ctx, "CreateCompany error = %v", err)
		return nil, err
	}
	return resp, nil
}

func Login(ctx context.Context, req *fabric_ebl.LoginReq) (resp *fabric_ebl.LoginResp, err error) {
	resp, err = client.Login(ctx, req)
	if err != nil {
		logger.CtxErrorf(ctx, "login error = %v", err)
		return nil, err
	}
	return resp, nil
}
