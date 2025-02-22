package common_user_rpc

import (
	"context"

	"ebl_api/common/wxl_cluster"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/wxl-server/idl_gen/kitex_gen/common_user"
	"github.com/wxl-server/idl_gen/kitex_gen/common_user/commonuser"
)

var client = wxl_cluster.NewClient(commonuser.NewClient, "common_user")

func Login(ctx context.Context, req *common_user.LoginReq) (resp *common_user.LoginResp, err error) {
	resp, err = client.Login(ctx, req)
	if err != nil {
		logger.CtxErrorf(ctx, "login error = %v", err)
		return nil, err
	}
	return resp, nil
}
func SignUp(ctx context.Context, req *common_user.SignUpReq) (resp *common_user.SignUpResp, err error) {
	resp, err = client.SignUp(ctx, req)
	if err != nil {
		logger.CtxErrorf(ctx, "SignUp error = %v", err)
		return nil, err
	}
	return resp, nil
}
