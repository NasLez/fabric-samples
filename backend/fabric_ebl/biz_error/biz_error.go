package biz_error

type BizError struct {
	code    int64
	message string
}

func (e *BizError) Error() string {
	return e.message
}

var (
	SignUpError          = newBizError(1000, "sign up failed, email has been used")
	UpdatePasswordError  = newBizError(1002, "update password failed, old password is same as new password")
	UpdatePasswordError2 = newBizError(1003, "update password failed, old password is wrong")
	LoginError           = newBizError(1004, "login failed, email is not exist")
	LoginError2          = newBizError(1005, "login failed, password is wrong")
	TokenError           = newBizError(1006, "token convert to claims failed")
	CreateCompanyError1  = newBizError(1007, "create company failed, company code has been used")
	CreateCompanyError2  = newBizError(1008, "create company failed, user email has been used")
	CreateCompanyError3  = newBizError(1009, "create company failed, user email is not exist")
	ParseTokenError      = newBizError(1010, "parse token failed")
	CompanyIDNotMatch    = newBizError(1011, "company id not match")
	StatusNotMatch       = newBizError(1012, "status not match")
	UserTypeNotMatch     = newBizError(1013, "user type not match")
)

func newBizError(code int64, message string) *BizError {
	return &BizError{
		code:    code,
		message: message,
	}
}
