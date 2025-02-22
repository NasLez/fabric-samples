package main

import (
	"context"
	"fabric_ipfs/service"
	fabric_ipfs "github.com/wxl-server/idl_gen/kitex_gen/fabric_ipfs"
	"go.uber.org/dig"
)

var handler fabric_ipfs.FabricIpfs

// FabricEblImpl implements the last service interface defined in the IDL.

// Handler implements the last service interface defined in the IDL.
type Handler struct {
	p Param
}

func (h Handler) CreateEblDocx(ctx context.Context, req *fabric_ipfs.CreateEblDocxReq) (r *fabric_ipfs.CreateEblDocxResp, err error) {
	return h.p.FabricEblService.CreateEblDocx(ctx, req)
}

type Param struct {
	dig.In
	FabricEblService service.FabricIpfsService
}

func NewHandler(p Param) {
	handler = &Handler{
		p: p,
	}
}
