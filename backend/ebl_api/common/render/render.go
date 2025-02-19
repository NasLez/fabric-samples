package render

import "ebl_api/common/json"

func Render(object any) string {
	return json.MarshalWithoutError[string](object)
}
