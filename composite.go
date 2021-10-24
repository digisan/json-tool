package jsontool

import (
	"log"

	"github.com/digisan/gotk/slice/ts"
	"github.com/tidwall/sjson"
)

// func Composite(m map[string]interface{}, filter func(path string) bool) string {
// 	jsonbytes, _ := sjson.SetBytes([]byte(""), "", "") // empty json doc to reinflate with tuples
// 	for path, value := range m {
// 		if (filter == nil) || (filter != nil && filter(path)) {
// 			jsonbytes, _ = sjson.SetBytes(jsonbytes, path, value)
// 		}
// 	}
// 	return string(jsonbytes)
// }

func Composite(m map[string]interface{}, fm func(path string, value interface{}) (p string, v interface{}, raw bool)) string {
	js, _ := sjson.Set("", "", "") // empty json doc to reinflate with tuples
	for path, value := range m {
		if fm == nil {
			js, _ = sjson.Set(js, path, value)
		} else {
			if p, v, raw := fm(path, value); p != "" && !raw {
				js, _ = sjson.Set(js, p, v)
			} else if p != "" && raw {
				js, _ = sjson.SetRaw(js, p, v.(string))
			}
		}
	}
	return js
}

func Composite2(m map[string]interface{}, fm func(path string, value interface{}) (p []string, v []interface{})) string {
	js, _ := sjson.Set("", "", "") // empty json doc to reinflate with tuples
	for path, value := range m {
		if fm == nil {
			js, _ = sjson.Set(js, path, value)
		} else {
			ps, vs := fm(path, value)
			if len(ps) != len(vs) {
				log.Fatalln("Composite2 [fm] return error")
			}
			for i, p := range ps {
				if p != "" {
					js, _ = sjson.Set(js, p, vs[i])
				}
			}
		}
	}
	return js
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
