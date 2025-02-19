package converter

import (
	"fabric_ebl/domain"
	"fabric_ebl/sal/dao/generator/model"
	"github.com/wxl-server/idl_gen/kitex_gen/fabric_ebl"
)

func CreateCompanyReqDTO2DO(req *fabric_ebl.CreateCompanyReq) *domain.CompanyDO {
	return &domain.CompanyDO{
		Name:    req.CompanyName,
		Code:    req.CompanyCode,
		Type:    int64(req.CompanyType),
		Balance: 500000,
	}
}

func CreateCompanyReqDTO2UserDO(req *fabric_ebl.CreateCompanyReq, companyId int64) *domain.UserDO {
	return &domain.UserDO{
		Email:     req.AdminEmail,
		Password:  req.AdminPassword,
		Name:      req.AdminName,
		Type:      1,
		Identity:  req.CompanyName + "." + req.AdminName,
		CompanyID: companyId,
	}
}

func CompanyDO2PO(do *domain.CompanyDO) *model.CompanyPO {
	return &model.CompanyPO{
		ID:        do.ID,
		Name:      do.Name,
		Code:      do.Code,
		Type:      do.Type,
		Seal:      do.Seal,
		Balance:   do.Balance,
		Extra:     do.Extra,
		CreatedAt: do.CreatedAt,
		UpdatedAt: do.UpdatedAt,
	}
}

func CompanyPO2DO(po *model.CompanyPO) *domain.CompanyDO {
	return &domain.CompanyDO{
		ID:        po.ID,
		Name:      po.Name,
		Code:      po.Code,
		Type:      po.Type,
		Seal:      po.Seal,
		Balance:   po.Balance,
		Extra:     po.Extra,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,
	}
}
