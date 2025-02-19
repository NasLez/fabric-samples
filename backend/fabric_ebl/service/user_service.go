package service

import (
	"context"
	"fabric_ebl/biz_error"
	"fabric_ebl/domain/converter"
	"fabric_ebl/repo"
	"fabric_ebl/sal/jwt"
	"github.com/wxl-server/idl_gen/kitex_gen/fabric_ebl"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/wxl-server/idl_gen/kitex_gen/common_user"
	"go.uber.org/dig"
)

type FabricEblService interface {
	SignUp(ctx context.Context, req *common_user.SignUpReq) (resp *common_user.SignUpResp, err error)
	UpdatePassword(ctx context.Context, req *common_user.UpdatePasswordReq) (*common_user.UpdatePasswordResp, error)
	Login(ctx context.Context, req *common_user.LoginReq) (resp *common_user.LoginResp, err error)
	CreateCompany(ctx context.Context, req *fabric_ebl.CreateCompanyReq) (resp *fabric_ebl.CreateCompanyResp, err error)
}

type Param struct {
	dig.In
	FabricEblRepo repo.FabricEblRepo
}

type FabricEblServiceImpl struct {
	p Param
}

func NewUserService(p Param) FabricEblService {
	return &FabricEblServiceImpl{
		p: p,
	}
}

func (u FabricEblServiceImpl) SignUp(ctx context.Context, req *common_user.SignUpReq) (resp *common_user.SignUpResp, err error) {
	count, err := u.p.FabricEblRepo.CountUser(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		logger.CtxErrorf(ctx, "sign up failed, email has been used, email = %s", req.Email)
		return nil, biz_error.SignUpError
	}
	id, err := u.p.FabricEblRepo.CreateUser(ctx, converter.SignUpReqDTO2DO(req))
	if err != nil {
		logger.CtxErrorf(ctx, "CreateUser failed, err = %v", err)
		return nil, err
	}
	return &common_user.SignUpResp{
		Id: id,
	}, nil
}

func (u FabricEblServiceImpl) UpdatePassword(ctx context.Context, req *common_user.UpdatePasswordReq) (*common_user.UpdatePasswordResp, error) {
	if req.OldPassword == req.Password {
		logger.CtxErrorf(ctx, "UpdatePassword failed, err = %v", biz_error.UpdatePasswordError)
		return nil, biz_error.UpdatePasswordError
	}
	do, err := u.p.FabricEblRepo.QueryUser(ctx, req.Email)
	if err != nil {
		logger.CtxErrorf(ctx, "QueryUser failed, err = %v", err)
		return nil, err
	}
	if do.Password != req.OldPassword {
		logger.CtxErrorf(ctx, "UpdatePassword failed, err = %v", biz_error.UpdatePasswordError2)
		return nil, biz_error.UpdatePasswordError2
	}
	do.Password = req.Password
	err = u.p.FabricEblRepo.UpdatePassword(ctx, do)
	if err != nil {
		logger.CtxErrorf(ctx, "UpdatePassword failed, err = %v", err)
		return nil, err
	}
	return &common_user.UpdatePasswordResp{}, nil
}

func (u FabricEblServiceImpl) Login(ctx context.Context, req *common_user.LoginReq) (resp *common_user.LoginResp, err error) {
	do, err := u.p.FabricEblRepo.QueryUser(ctx, req.Email)
	if err != nil {
		logger.CtxErrorf(ctx, "QueryUser failed, err = %v", err)
		return nil, biz_error.LoginError
	}
	if do.Password != req.Password {
		logger.CtxErrorf(ctx, "login failed, password is wrong")
		return nil, biz_error.LoginError2
	}
	token, err := jwt.GenerateToken(ctx, do.ID)
	if err != nil {
		logger.CtxErrorf(ctx, "login failed, err = %v", err)
		return nil, biz_error.LoginError
	}
	return &common_user.LoginResp{
		Token: token,
	}, nil
}

func (u FabricEblServiceImpl) CreateCompany(ctx context.Context, req *fabric_ebl.CreateCompanyReq) (resp *fabric_ebl.CreateCompanyResp, err error) {
	{
		count, err := u.p.FabricEblRepo.CountCompanyByCode(ctx, req.CompanyCode)
		if err != nil {
			logger.CtxErrorf(ctx, "QueryCompanyByCode failed, err = %v", err)
			return nil, err
		}
		if count > 0 {
			logger.CtxErrorf(ctx, "CreateCompany failed, company code has been created, companyCode = %s", req.CompanyCode)
			return nil, biz_error.CreateCompanyError1
		}
	}
	{
		count, err := u.p.FabricEblRepo.CountUser(ctx, req.AdminEmail)
		if err != nil {
			logger.CtxErrorf(ctx, "QueryUser failed, err = %v", err)
			return nil, err
		}
		if count > 0 {
			logger.CtxErrorf(ctx, "CreateCompany failed, user has been created, email = %s", req.AdminEmail)
			return nil, biz_error.CreateCompanyError2
		}
	}
	companyId, err := u.p.FabricEblRepo.CreateCompany(ctx, converter.CreateCompanyReqDTO2DO(req))
	if err != nil {
		logger.CtxErrorf(ctx, "CreateCompany failed, err = %v", err)
		return nil, err
	}
	userId, err := u.p.FabricEblRepo.CreateUser(ctx, converter.CreateCompanyReqDTO2UserDO(req, companyId))
	if err != nil {
		logger.CtxErrorf(ctx, "CreateUser failed, err = %v", err)
		return nil, err
	}
	return &fabric_ebl.CreateCompanyResp{
		Id: userId,
	}, nil
}
