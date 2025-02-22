package domain

type UserDO struct {
	ID        int64
	Email     string
	Password  string
	Name      string
	Type      int64
	CompanyID int64
	Identity  string
	Extra     *string
	CreatedAt int64
	UpdatedAt int64
}
