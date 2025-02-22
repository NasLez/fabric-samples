package converter

import (
	"fabric_ipfs/domain"
	"fabric_ipfs/sal/dao/generator/model"

	"github.com/wxl-server/idl_gen/kitex_gen/common_user"
)

func SignUpReqDTO2DO(req *common_user.SignUpReq) *domain.UserDO {
	return &domain.UserDO{
		Email:    req.Email,
		Password: req.Password,
	}
}

func UserDO2PO(do *domain.UserDO) *model.UserPO {
	return &model.UserPO{
		ID:        do.ID,
		Email:     do.Email,
		Password:  do.Password,
		Name:      do.Name,
		Type:      do.Type,
		CompanyID: do.CompanyID,
		Identity:  &do.Identity,
		Extra:     do.Extra,
		CreatedAt: do.CreatedAt,
		UpdatedAt: do.UpdatedAt,
	}
}

func UserPO2DO(po *model.UserPO) *domain.UserDO {
	return &domain.UserDO{
		ID:        po.ID,
		Email:     po.Email,
		Password:  po.Password,
		Name:      po.Name,
		Type:      po.Type,
		CompanyID: po.CompanyID,
		Identity:  *po.Identity,
		Extra:     po.Extra,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,
	}
}
