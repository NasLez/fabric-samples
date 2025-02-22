package Status

import "fabric_ipfs/common/gptr"

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
)
