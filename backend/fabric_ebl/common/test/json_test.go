package test

import (
	"fabric_ebl/common/json"
	"log"
	"testing"
)

func TestJson(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	testStruct := TestStruct{
		Name: "test",
		Age:  18,
	}
	jsonStr := json.MarshalWithoutError[string](testStruct)
	log.Print(jsonStr)
	testStruct2 := json.UnmarshalWithoutError[TestStruct](jsonStr)
	log.Print(json.MarshalWithoutError[string](testStruct2))
	testStruct3 := json.UnmarshalWithoutError[*TestStruct](jsonStr)
	log.Print(json.MarshalWithoutError[string](testStruct3))

	type TestStruct2 struct {
		Named string `json:"named"`
		Ages  int    `json:"ages"`
	}
	testStruct4 := json.UnmarshalWithoutError[TestStruct2](jsonStr)
	log.Print(json.MarshalWithoutError[string](testStruct4))
}
