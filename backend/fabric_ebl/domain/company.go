package domain

type CompanyDO struct {
	ID        int64
	Name      string
	Code      string
	Seal      *string
	Balance   float64
	Type      int64
	Extra     *string
	CreatedAt int64
	UpdatedAt int64
}

type QueryCompanyDo struct {
	ID   int64
	Name string
	Code string
	Type int64
}
