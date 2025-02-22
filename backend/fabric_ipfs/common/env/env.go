package env

import "runtime"

type Env = string

const (
	Prod Env = "prod"
	Boe  Env = "boe"
)

func GetEnv() Env {
	sysType := runtime.GOOS

	//if sysType == "linux" {
	//	return Prod
	//}
	if sysType == "linux" {
		return Boe
	}

	if sysType == "windows" {
		return Boe
	}
	return Boe
}

func IsProd() bool {
	return GetEnv() == Prod
}

func IsBoe() bool {
	return GetEnv() == Boe
}
