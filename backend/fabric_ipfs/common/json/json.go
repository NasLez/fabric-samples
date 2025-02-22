package json

import gojson "encoding/json"

type strTypes interface {
	[]byte | string
}

func MarshalWithoutError[T strTypes](object any) T {
	marshal, _ := gojson.Marshal(object)
	return T(marshal)
}

func UnmarshalWithoutError[T any, M strTypes](jsonStr M) (object T) {
	_ = gojson.Unmarshal([]byte(jsonStr), &object)
	return
}
