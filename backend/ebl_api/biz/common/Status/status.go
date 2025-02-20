package Status

import "github.com/wxl-server/common/gptr"

type Status struct {
	code    *int64
	message *string
}

func (s *Status) Code() *int64 {
	return s.code
}

func (s *Status) Message() *string {
	return s.message
}

var (
	Success = &Status{
		code:    gptr.Of(int64(0)),
		message: gptr.Of("success"),
	}
	RequestParamsInvalid = &Status{
		code:    gptr.Of(int64(1)),
		message: gptr.Of("request params invalid"),
	}
	LoginError = &Status{
		code:    gptr.Of(int64(1)),
		message: gptr.Of("email or password is wrong"),
	}
	CreateCompanyError = &Status{
		code:    gptr.Of(int64(1)),
		message: gptr.Of("create company error"),
	}
	GetUserInfoError = &Status{
		code:    gptr.Of(int64(1)),
		message: gptr.Of("get user info error"),
	}
	GetCompanyAllListError = &Status{
		code:    gptr.Of(int64(1)),
		message: gptr.Of("get company all list error"),
	}
	CreateEblError = &Status{
		code:    gptr.Of(int64(1)),
		message: gptr.Of("create ebl error"),
	}
)
