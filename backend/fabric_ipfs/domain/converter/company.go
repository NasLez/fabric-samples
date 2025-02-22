package converter

import (
	"fabric_ipfs/domain"
	"fabric_ipfs/sal/dao/generator/model"
)

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

func CompanyPOs2DOs(pos []*model.CompanyPO) (dos []domain.QueryCompanyDo) {
	dos = make([]domain.QueryCompanyDo, 0, len(pos))
	for po := range pos {
		dos = append(dos, CompanyPo2QueryCompanyDo(pos[po]))
	}
	return dos
}

func CompanyPo2QueryCompanyDo(po *model.CompanyPO) domain.QueryCompanyDo {
	return domain.QueryCompanyDo{
		ID:   po.ID,
		Name: po.Name,
		Code: po.Code,
		Type: po.Type,
	}
}
