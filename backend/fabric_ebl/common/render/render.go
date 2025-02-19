package render

import "fabric_ebl/common/json"

func Render(object any) string {
	return json.MarshalWithoutError[string](object)
}
