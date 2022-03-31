package jsontool

import (
	"log"
	"math"
	"regexp"

	. "github.com/digisan/go-generics/v2"
	"github.com/tidwall/sjson"
)

func GetStrVal(value any) string {
	val := ""
	switch v := value.(type) {
	case string:
		val = v
	}
	return val
}

func GetIntVal(value any) int {
	val := math.MinInt + 1
	switch v := value.(type) {
	case int:
		val = v
	}
	return val
}

var (
	// mRule  = make(map[string]func(path string, value any) (ok bool, ps []string, vs []any))
	// mFocus = make(map[string]*regexp.Regexp)
	focus = []*regexp.Regexp{}
	rules = []func(path string, value any) (ok bool, ps []string, vs []any){}
)

// func RegisterRule(name string, regexpStr string, f func(path string, value any) (ok bool, ps []string, vs []any)) {
// 	mRule[name] = f
// 	mFocus[name] = regexp.MustCompile(regexpStr)
// }

func RegisterRule(regexpStr string, f func(path string, value any) (ok bool, ps []string, vs []any)) {
	focus = append(focus, regexp.MustCompile(regexpStr))
	rules = append(rules, f)
}

func TransformUnderFirstRule(mData map[string]any, data []byte) string {

	var err error
	if mData == nil {
		if mData, err = Flatten(data); err != nil {
			log.Fatalln(err)
		}
	}

	js, _ := sjson.Set("", "", "") // empty json doc to reinflate with tuples
NEXT_PATH:
	for path, value := range mData {
		for iR, rule := range rules {
			if focus[iR].MatchString(path) { // only process focused path
				if ok, ps, vs := rule(path, value); ok { // only process 'ok' condition
					if len(ps) != len(vs) {
						log.Fatalln("Transform [rule] return error")
					}
					for i, p := range ps {
						if p != "" { // non empty path, modify result
							if js, err = sjson.Set(js, p, vs[i]); err != nil {
								log.Fatalln(err)
							}
						}
						// empty path ("") => delete this path
					}
					continue NEXT_PATH
				}
			}
		}
		// no ruled, keep original path-value, ignore further processing
		if js, err = sjson.Set(js, path, value); err != nil {
			log.Fatalln(err)
		}
	}
	return js
}

// func TransformUnderLastRule(data []byte) string {
// 	mData, err := Flatten(data)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	js, _ := sjson.Set("", "", "") // empty json doc to reinflate with tuples
// 	for path, value := range mData {
// 		ruled := false
// 		for iR, rule := range rules {
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
// 					ruled = true
// 				}
// 			}
// 		}
// 		// no ruled, keep original path-value, ignore further processing
// 		if !ruled {
// 			js, _ = sjson.Set(js, path, value)
// 		}
// 	}
// 	return js
// }

func TransformUnderAllRules(mData map[string]any, data []byte) string {

	var err error
	if mData == nil {
		if mData, err = Flatten(data); err != nil {
			log.Fatalln(err)
		}
	}

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
							if js, err = sjson.Set(js, p, vs[i]); err != nil {
								log.Fatalln(err)
							}
						}
						// empty path ("") => delete original path
					}
				} else { // no ok, keep original path-value, ignore further processing
					if js, err = sjson.Set(js, path, value); err != nil {
						log.Fatalln(err)
					}
				}
			} else { // no focused, keep original path-value, ignore further processing
				if js, err = sjson.Set(js, path, value); err != nil {
					log.Fatalln(err)
				}
			}
		}
		if mData, err = FlattenStr(js); err != nil {
			log.Fatalln(err)
		}
	}
	return js
}

// func Composite(m map[string]any, filter func(path string) bool) string {
// 	jsonbytes, _ := sjson.SetBytes([]byte(""), "", "") // empty json doc to reinflate with tuples
// 	for path, value := range m {
// 		if (filter == nil) || (filter != nil && filter(path)) {
// 			jsonbytes, _ = sjson.SetBytes(jsonbytes, path, value)
// 		}
// 	}
// 	return string(jsonbytes)
// }

func Composite(m map[string]any, fm func(path string, value any) (p string, v any, raw bool)) string {
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

func Composite2(m map[string]any, fm func(path string, value any) (p []string, v []any)) string {
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

func CompositeExcl(m map[string]any, exclPaths ...string) string {
	jsonbytes, _ := sjson.SetBytes([]byte(""), "", "") // empty json doc to reinflate with tuples
	for path, value := range m {
		if exclPaths != nil && In(path, exclPaths...) {
			continue
		}
		jsonbytes, _ = sjson.SetBytes(jsonbytes, path, value)
	}
	return string(jsonbytes)
}

func CompositeIncl(m map[string]any, inclPaths ...string) string {
	jsonbytes, _ := sjson.SetBytes([]byte(""), "", "") // empty json doc to reinflate with tuples
	for path, value := range m {
		if inclPaths != nil && In(path, inclPaths...) {
			jsonbytes, _ = sjson.SetBytes(jsonbytes, path, value)
		}
	}
	return string(jsonbytes)
}
