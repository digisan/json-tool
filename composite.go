package jsontool

import (
	"github.com/tidwall/sjson"
)

func Composite(m map[string]interface{}) string {
	jsonbytes, _ := sjson.SetBytes([]byte(""), "", "") // empty json doc to reinflate with tuples
	for path, value := range m {
		jsonbytes, _ = sjson.SetBytes(jsonbytes, path, value)
	}
	return string(jsonbytes)
}
