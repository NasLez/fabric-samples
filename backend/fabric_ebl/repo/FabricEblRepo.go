package repo

import (
	"context"
	"fabric_ebl/domain"
	"fabric_ebl/domain/converter"
	"fabric_ebl/sal/dao/generator/query"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/wxl-server/common/id_gen"

	"go.uber.org/dig"
)

type FabricEblRepo interface {
	CountUser(ctx context.Context, email string) (count int64, err error)
	CreateUser(ctx context.Context, do *domain.UserDO) (id int64, err error)
	QueryUser(ctx context.Context, email string) (do *domain.UserDO, err error)
	UpdatePassword(ctx context.Context, do *domain.UserDO) (err error)
	QueryCompanyByCode(ctx context.Context, code string) (do *domain.CompanyDO, err error)
	CountCompanyByCode(ctx context.Context, code string) (count int64, err error)
	CreateCompany(ctx context.Context, do *domain.CompanyDO) (id int64, err error)
	QueryUserById(ctx context.Context, id int64) (do *domain.UserDO, err error)
	QueryCompanyById(ctx context.Context, id int64) (do *domain.CompanyDO, err error)
	QueryCompanyAll(ctx context.Context) (do []domain.QueryCompanyDo, err error)
	UpdateCompanySeal(ctx context.Context, do *domain.CompanyDO) (err error)
}

type Param struct {
	dig.In
}

type FabricEblRepoImpl struct {
	p Param
}

func (u FabricEblRepoImpl) UpdateCompanySeal(ctx context.Context, do *domain.CompanyDO) (err error) {
	po := query.Q.CompanyPO
	_, err = po.WithContext(ctx).Where(po.ID.Eq(do.ID)).Update(po.Seal, do.Seal)
	if err != nil {
		logger.CtxErrorf(ctx, "UpdateCompanySeal failed, err = %v", err)
		return err
	}
	return nil
}

func (u FabricEblRepoImpl) QueryCompanyAll(ctx context.Context) (do []domain.QueryCompanyDo, err error) {
	po := query.Q.CompanyPO
	condition := po.WithContext(ctx)
	companyPOS, err := condition.Find()
	if err != nil {
		logger.CtxErrorf(ctx, "condition.Find failed, err = %v", err)
		return nil, err
	}
	return converter.CompanyPOs2DOs(companyPOS), nil
}

func (u FabricEblRepoImpl) QueryCompanyById(ctx context.Context, id int64) (do *domain.CompanyDO, err error) {
	po := query.Q.CompanyPO
	companyPO, err := po.WithContext(ctx).Where(po.ID.Eq(id)).First()
	if err != nil {
		logger.CtxErrorf(ctx, "QueryCompanyById failed, err", err)
		return nil, err
	}
	return converter.CompanyPO2DO(companyPO), nil
}

func (u FabricEblRepoImpl) QueryUserById(ctx context.Context, id int64) (do *domain.UserDO, err error) {
	po := query.Q.UserPO
	userPO, err := po.WithContext(ctx).Where(po.ID.Eq(id)).First()
	if err != nil {
		logger.CtxErrorf(ctx, "QueryUserById failed, err = %v", err)
		return nil, err
	}
	return converter.UserPO2DO(userPO), nil
}

func (u FabricEblRepoImpl) CreateCompany(ctx context.Context, do *domain.CompanyDO) (id int64, err error) {
	po := query.Q.CompanyPO
	do.ID, err = id_gen.NextID()
	if err != nil {
		logger.CtxErrorf(ctx, "generate id failed, err = %v", err)
		return 0, err
	}
	err = po.WithContext(ctx).Create(converter.CompanyDO2PO(do))
	if err != nil {
		logger.CtxErrorf(ctx, "CreateCompany failed, err = %v", err)
		return 0, err
	}
	return do.ID, nil
}

func NewFabricEblRepo(p Param) FabricEblRepo {
	return &FabricEblRepoImpl{
		p: p,
	}
}

func (u FabricEblRepoImpl) UpdatePassword(ctx context.Context, do *domain.UserDO) (err error) {
	po := query.Q.UserPO
	_, err = po.WithContext(ctx).Where(po.ID.Eq(do.ID)).Update(po.Password, do.Password)
	if err != nil {
		logger.CtxErrorf(ctx, "UpdatePassword failed, err = %v", err)
		return err
	}
	return nil
}

func (u FabricEblRepoImpl) QueryUser(ctx context.Context, email string) (do *domain.UserDO, err error) {
	po := query.Q.UserPO
	userPO, err := po.WithContext(ctx).Where(po.Email.Eq(email)).First()
	if err != nil {
		logger.CtxErrorf(ctx, "QueryUser failed, err = %v", err)
		return nil, err
	}
	return converter.UserPO2DO(userPO), nil
}

func (u FabricEblRepoImpl) CountUser(ctx context.Context, email string) (count int64, err error) {
	po := query.Q.UserPO
	count, err = po.WithContext(ctx).Where(po.Email.Eq(email)).Count()
	if err != nil {
		logger.CtxErrorf(ctx, "CountUser failed, err = %v", err)
		return 0, err
	}
	return count, nil
}

func (u FabricEblRepoImpl) CreateUser(ctx context.Context, do *domain.UserDO) (id int64, err error) {
	po := query.Q.UserPO
	do.ID, err = id_gen.NextID()
	if err != nil {
		logger.CtxErrorf(ctx, "generate id failed, err = %v", err)
		return 0, err
	}
	err = po.WithContext(ctx).Create(converter.UserDO2PO(do))
	if err != nil {
		logger.CtxErrorf(ctx, "CreateUser failed, err = %v", err)
		return 0, err
	}
	return do.ID, nil
}

func (u FabricEblRepoImpl) QueryCompanyByCode(ctx context.Context, code string) (do *domain.CompanyDO, err error) {
	po := query.Q.CompanyPO
	companyPO, err := po.WithContext(ctx).Where(po.Code.Eq(code)).First()
	if err != nil {
		logger.CtxErrorf(ctx, "QueryCompanyByCode failed, err = %v", err)
		return nil, err
	}
	return converter.CompanyPO2DO(companyPO), nil
}

func (u FabricEblRepoImpl) CountCompanyByCode(ctx context.Context, code string) (count int64, err error) {
	po := query.Q.CompanyPO
	count, err = po.WithContext(ctx).Where(po.Code.Eq(code)).Count()
	if err != nil {
		logger.CtxErrorf(ctx, "CountCompanyByCode failed, err = %v", err)
		return 0, err
	}
	return count, nil
}
