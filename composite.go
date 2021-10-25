package jsontool

import (
	"log"
	"math"
	"regexp"

	"github.com/digisan/gotk/slice/ts"
	"github.com/tidwall/sjson"
)

func GetStrVal(value interface{}) string {
	val := ""
	switch v := value.(type) {
	case string:
		val = v
	}
	return val
}

func GetIntVal(value interface{}) int {
	val := math.MinInt + 1
	switch v := value.(type) {
	case int:
		val = v
	}
	return val
}

var (
	// mRule  = make(map[string]func(path string, value interface{}) (ok bool, ps []string, vs []interface{}))
	// mFocus = make(map[string]*regexp.Regexp)
	focus = []*regexp.Regexp{}
	rules = []func(path string, value interface{}) (ok bool, ps []string, vs []interface{}){}
)

// func RegisterRule(name string, regexpStr string, f func(path string, value interface{}) (ok bool, ps []string, vs []interface{})) {
// 	mRule[name] = f
// 	mFocus[name] = regexp.MustCompile(regexpStr)
// }

func RegisterRule(regexpStr string, f func(path string, value interface{}) (ok bool, ps []string, vs []interface{})) {
	focus = append(focus, regexp.MustCompile(regexpStr))
	rules = append(rules, f)
}

func Transform(data []byte) string {
	mData, err := Flatten(data)
	if err != nil {
		log.Fatalln(err)
	}

	// for path, value := range mData {
	// 	if len(rules) == 0 {
	// 		js, _ = sjson.Set(js, path, value)
	// 		continue
	// 	}
	// 	for iR, rule := range rules {
	// 		if rule == nil { // no rule, ignore further processing
	// 			js, _ = sjson.Set(js, path, value)
	// 		} else {
	// 			if focus[iR].MatchString(path) { // only process focused path
	// 				if ok, ps, vs := rule(path, value); ok { // only process 'ok' condition
	// 					if len(ps) != len(vs) {
	// 						log.Fatalln("Transform [rule] return error")
	// 					}
	// 					for i, p := range ps {
	// 						if p != "" { // non empty path, modify result
	// 							js, _ = sjson.Set(js, p, vs[i])
	// 						}
	// 						// empty path ("") => delete this path
	// 					}
	// 				} else { // no ok, keep original path-value, ignore further processing
	// 					js, _ = sjson.Set(js, path, value)
	// 				}
	// 			} else { // no focused, keep original path-value, ignore further processing
	// 				js, _ = sjson.Set(js, path, value)
	// 			}
	// 		}
	// 	}
	// }

	js, _ := sjson.Set("", "", "") // empty json doc to reinflate with tuples
	for iR, rule := range rules {
		js, _ = sjson.Set("", "", "")
		for path, value := range mData {
			if focus[iR].MatchString(path) { // only process focused path
				if ok, ps, vs := rule(path, value); ok { // only process 'ok' condition
					if len(ps) != len(vs) {
						log.Fatalln("Transform [rule] return error")
					}
					for i, p := range ps {
						if p != "" { // non empty path, modify result
							js, _ = sjson.Set(js, p, vs[i])
						}
						// empty path ("") => delete original path
					}
				} else { // no ok, keep original path-value, ignore further processing
					js, _ = sjson.Set(js, path, value)
				}
			} else { // no focused, keep original path-value, ignore further processing
				js, _ = sjson.Set(js, path, value)
			}
		}
		if mData, err = FlattenStr(js); err != nil {
			log.Fatalln(err)
		}
	}

	return js
}

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
