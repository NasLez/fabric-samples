package test

import (
	"fabric_ebl/common/env"
	"log"
	"testing"
)

func TestEnv(t *testing.T) {
	envStr := env.GetEnv()
	log.Print(envStr)
	log.Print(env.IsProd())
	log.Print(env.IsBoe())
}
