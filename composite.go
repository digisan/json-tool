package jsontool

import (
	"github.com/digisan/gotk/slice/ts"
	"github.com/tidwall/sjson"
)

func Composite(m map[string]interface{}, filter func(path string) bool) string {
	jsonbytes, _ := sjson.SetBytes([]byte(""), "", "") // empty json doc to reinflate with tuples
	for path, value := range m {
		if filter != nil && filter(path) {
			jsonbytes, _ = sjson.SetBytes(jsonbytes, path, value)
		}
	}
	return string(jsonbytes)
}

func CompositeExcl(m map[string]interface{}, exclPaths ...string) string {
	jsonbytes, _ := sjson.SetBytes([]byte(""), "", "") // empty json doc to reinflate with tuples
	for path, value := range m {
		if exclPaths != nil && ts.In(path, exclPaths...) {
			continue
		}
		jsonbytes, _ = sjson.SetBytes(jsonbytes, path, value)
	}
	return string(jsonbytes)
}

func CompositeIncl(m map[string]interface{}, inclPaths ...string) string {
	jsonbytes, _ := sjson.SetBytes([]byte(""), "", "") // empty json doc to reinflate with tuples
	for path, value := range m {
		if inclPaths != nil && ts.In(path, inclPaths...) {
			jsonbytes, _ = sjson.SetBytes(jsonbytes, path, value)
		}
	}
	return string(jsonbytes)
}
