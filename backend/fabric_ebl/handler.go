package main

import (
	"context"
	"fabric_ebl/service"
	"github.com/wxl-server/idl_gen/kitex_gen/fabric_ebl"
	"go.uber.org/dig"
)

var handler fabric_ebl.FabricEbl

// FabricEblImpl implements the last service interface defined in the IDL.

// Handler implements the last service interface defined in the IDL.
type Handler struct {
	p Param
}

func (s *Handler) CreateEbl(ctx context.Context, req *fabric_ebl.CreateEblReq) (r *fabric_ebl.CreateEblResp, err error) {
	return s.p.FabricEblService.CreateEbl(ctx, req)
}

func (s *Handler) GetCompanyAllList(ctx context.Context, req *fabric_ebl.GetCompanyAllListReq) (r *fabric_ebl.GetCompanyAllListResp, err error) {
	return s.p.FabricEblService.GetCompanyAllList(ctx, req)
}

func (s *Handler) GetUserInfo(ctx context.Context, req *fabric_ebl.GetUserInfoReq) (r *fabric_ebl.GetUserInfoResp, err error) {
	return s.p.FabricEblService.GetUserInfo(ctx, req)
}

func (s *Handler) Login(ctx context.Context, req *fabric_ebl.LoginReq) (r *fabric_ebl.LoginResp, err error) {
	return s.p.FabricEblService.Login(ctx, req)
}

type Param struct {
	dig.In
	FabricEblService service.FabricEblService
}

func NewHandler(p Param) {
	handler = &Handler{
		p: p,
	}
}

func (s *Handler) CreateCompany(ctx context.Context, req *fabric_ebl.CreateCompanyReq) (r *fabric_ebl.CreateCompanyResp, err error) {
	return s.p.FabricEblService.CreateCompany(ctx, req)
}
