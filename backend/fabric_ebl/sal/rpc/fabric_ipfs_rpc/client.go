package fabric_ipfs_rpc

import (
	"context"
	"github.com/wxl-server/idl_gen/kitex_gen/fabric_ipfs"
	"github.com/wxl-server/idl_gen/kitex_gen/fabric_ipfs/fabricipfs"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/wxl-server/common/wxl_cluster"
)

var client = wxl_cluster.NewClient(fabricipfs.NewClient, "fabric_ipfs")

func CreateEblDocx(ctx context.Context, req *fabric_ipfs.CreateEblDocxReq) (resp *fabric_ipfs.CreateEblDocxResp, err error) {
	resp, err = client.CreateEblDocx(ctx, req)
	if err != nil {
		logger.CtxErrorf(ctx, "CreateEblDocx error = %v", err)
		return nil, err
	}
	return resp, nil
}
