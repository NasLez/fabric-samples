package domain

// EblFilter 表示电子提单过滤器的结构体
type EblFilter struct {
	// EblNo 电子提单编号
	EblNo string `json:"eblNo"`
	// OriginCompanyID 起始公司 ID
	OriginCompanyID string `json:"originCompanyID"`
	// ShipperCompanyID 发货人公司 ID
	ShipperCompanyID string `json:"shipperCompanyID"`
	// ConsigneeCompanyID 收货人公司 ID
	ConsigneeCompanyID string `json:"consigneeCompanyID"`
	// NotifyPartyCompanyID 通知方公司 ID
	NotifyPartyCompanyID string `json:"notifyPartyCompanyID"`
	// PlaceOfReceipt 收货地点
	PlaceOfReceipt string `json:"placeOfReceipt"`
	// OceanVessel 远洋船舶
	OceanVessel string `json:"oceanVessel"`
	// PortOfLoading 装货港
	PortOfLoading string `json:"portOfLoading"`
	// PortOfDescharge 卸货港
	PortOfDescharge string `json:"portOfDescharge"`
	// PlaceOfDestination 目的地
	PlaceOfDestination string `json:"placeOfDestination"`
	// PlaceOfDelivery 交货地点
	PlaceOfDelivery string `json:"placeOfDelivery"`
	// ShippingMarkes 唛头
	ShippingMarkes string `json:"shippingMarkes"`
	// QuantityOfPackages 包装数量
	QuantityOfPackages float64 `json:"quantityOfPackages"`
	// KindOfPackagesGW 包装种类（毛重相关）
	KindOfPackagesGW string `json:"kindOfPackagesGW"`
	// KindOfPackagesM 包装种类（体积相关）
	KindOfPackagesM string `json:"kindOfPackagesM"`
	// DescriptionOfGoods 货物描述
	DescriptionOfGoods string `json:"descriptionOfGoods"`
	// GrossWeight 毛重
	GrossWeight float64 `json:"grossWeight"`
	// Measurement 体积
	Measurement float64 `json:"measurement"`
	// FreightAndCharges 运费和费用
	FreightAndCharges string `json:"freightAndCharges"`
	// PlaceOfIssue 签发地点
	PlaceOfIssue string `json:"placeOfIssue"`
	// DateOfIssue 签发日期
	DateOfIssue int64 `json:"dateOfIssue"`
	// DeliveryAgent 交付代理
	DeliveryAgent string `json:"deliveryAgent"`
	// ShippedOnBoard 装船日期
	ShippedOnBoard int64 `json:"shippedOnBoard"`
	// NumOfEBL 电子提单数量
	NumOfEBL int64 `json:"numOfEBL"`
	// DateOfIssueDeadline 签发截止日期
	DateOfIssueDeadline int64 `json:"dateOfIssueDeadline"`
	// Status 状态
	Status string `json:"status"`
	// TransferCompanyID 转运公司 ID
	TransferCompanyID string `json:"transferCompanyID"`
	// CompanyID 公司 ID
	CompanyID int64 `json:"companyID"`
}
