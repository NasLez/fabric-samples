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

func GetUserInfo(ctx context.Context, req *fabric_ebl.GetUserInfoReq) (resp *fabric_ebl.GetUserInfoResp, err error) {
	resp, err = client.GetUserInfo(ctx, req)
	if err != nil {
		logger.CtxErrorf(ctx, "GetUserInfo error = %v", err)
		return nil, err
	}
	return resp, nil
}

func GetCompanyAllList(ctx context.Context, rpc *fabric_ebl.GetCompanyAllListReq) (resp *fabric_ebl.GetCompanyAllListResp, err error) {
	resp, err = client.GetCompanyAllList(ctx, rpc)
	if err != nil {
		logger.CtxErrorf(ctx, "GetCompanyAllList error = %v", err)
		return nil, err
	}
	return resp, nil
}

func CreateEbl(ctx context.Context, rpc *fabric_ebl.CreateEblReq) (resp *fabric_ebl.CreateEblResp, err error) {
	resp, err = client.CreateEbl(ctx, rpc)
	if err != nil {
		logger.CtxErrorf(ctx, "CreateEbl error = %v", err)
		return nil, err
	}
	return resp, nil
}

func QueryAllEblList(ctx context.Context, rpc *fabric_ebl.QueryAllEblListReq) (resp *fabric_ebl.QueryAllEblListResp, err error) {
	resp, err = client.QueryAllEblList(ctx, rpc)
	if err != nil {
		logger.CtxErrorf(ctx, "QueryAllEblList error = %v", err)
		return nil, err
	}
	return resp, nil
}
