package render

import "fabric_ipfs/common/json"

func Render(object any) string {
	return json.MarshalWithoutError[string](object)
}
