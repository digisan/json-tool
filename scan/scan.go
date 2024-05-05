package scan

import (
	"fmt"
	"log"
	"os"
	"strings"

	. "github.com/digisan/go-generics"
	fd "github.com/digisan/gotk/file-dir"
	"github.com/digisan/gotk/strs"
	jt "github.com/digisan/json-tool"
)

// *** !!! MUST be JS formatted !!! *** //

const (
	JUNK   = "**"
	INDENT = "  "
)

const (
	OBJ_OPEN  = "{"
	OBJ_CLOSE = "}"
	OBJ_EMPTY = "{}"
	ARR_OPEN  = "["
	ARR_CLOSE = "]"
	ARR_EMPTY = "[]"
)

type LineType int

const (
	KV          LineType = iota // e.g. "code": null(,)
	KV_STR                      // e.g. "code": "root"(,)
	KV_OBJ_OPEN                 // e.g. "doc": {
	KV_ARR_OPEN                 // e.g. "children": [
	OBJ                         // e.g. { OR } OR {}(,)
	ARR                         // e.g. [ OR ] OR [](,)
	ELEM                        // e.g. true OR false OR null OR 1,2,3(,)
	ELEM_STR                    // e.g. "AC9RMDPS_FY"(,)
	UNKNOWN
)

func simpleFetch(line string) (lnType LineType, k string, v any, comma bool) {

	lnType = UNKNOWN
	comma = strings.HasSuffix(line, ",")                   // check if ',' at the end
	ln := strings.TrimSuffix(strings.TrimSpace(line), ",") // trim spaces & ended ','

	switch {

	// ** "key": value(raw string) ** //
	case ln[0] == '"' && strings.Contains(ln, `": `):

		kv := strings.Split(ln, `": `)
		k = kv[0][1:]
		vstr := kv[1]

		if vstr[0] == '"' && vstr[len(vstr)-1] == '"' { // string value

			v = vstr[1 : len(vstr)-1]
			lnType = KV_STR

		} else if In(vstr, OBJ_OPEN, OBJ_EMPTY) { // object open OR empty object value

			v = vstr
			lnType = KV_OBJ_OPEN

		} else if In(vstr, ARR_OPEN, ARR_EMPTY) { // array open OR empty array value

			v = vstr
			lnType = KV_ARR_OPEN

		} else { // other non-string value, like null, 1, 2, true etc.

			v = vstr
			lnType = KV
		}

		// ** No key; value: object open OR close OR empty, i,e { OR } OR {} ** //
	case In(ln, OBJ_OPEN, OBJ_CLOSE, OBJ_EMPTY):
		k = ""
		v = ln
		lnType = OBJ

		// ** No key; value: array open or close OR empty, i.e [ OR ] OR [] ** //
	case In(ln, ARR_OPEN, ARR_CLOSE, ARR_EMPTY):
		k = ""
		v = ln
		lnType = ARR

		// ** No key; value: string single array element line, e.g. "AC9RMDPS_FY" ** //
	case ln[0] == '"' && ln[len(ln)-1] == '"':
		k = ""
		v = ln[1 : len(ln)-1]
		lnType = ELEM_STR

		// ** No key; value: non-string single array element line, e.g. true/false/null OR 1/2/3 etc. ** //
	default:
		k = ""
		v = ln
		lnType = ELEM
	}
	return
}

type LineInfo struct {
	line   string
	ln     string
	key    string
	lnType LineType
	lvl    int
}

func fnSetCurrentKeyLevel() func(above, this, below LineInfo) string {

	mLvlKey := make(map[int]string)
	mPathIdx := make(map[string]int)

	clrAfter := func(lvl int) {
		for kLvl := range mLvlKey {
			if kLvl > lvl {
				delete(mLvlKey, kLvl)
			}
		}
	}

	clrThisAndAfter := func(lvl int) {
		for kLvl := range mLvlKey {
			if kLvl >= lvl {
				delete(mLvlKey, kLvl)
			}
		}
	}

	return func(above, this, below LineInfo) string {

		var (
			flagCloseObj = false
			flagCloseArr = false
			// flagEmptyObj = false
			// flagEmptyArr = false
		)

		if In(this.lnType, KV, KV_STR, KV_OBJ_OPEN, KV_ARR_OPEN) {
			mLvlKey[this.lvl] = this.key
			clrAfter(this.lvl)
		}

		if this.ln == OBJ_CLOSE {
			clrAfter(this.lvl)
			flagCloseObj = true
		}
		if this.ln == ARR_CLOSE {
			clrAfter(this.lvl)
			flagCloseArr = true
		}
		// if this.ln == OBJ_EMPTY {
		// 	clrAfter(this.lvl)
		// 	flagEmptyObj = true
		// }
		// if this.ln == ARR_EMPTY {
		// 	clrAfter(this.lvl)
		// 	flagEmptyArr = true
		// }

		if In(this.ln, OBJ_OPEN, OBJ_EMPTY) {

			clrThisAndAfter(this.lvl)
			_, values := MapToKVs(mLvlKey, func(ki, kj int) bool {
				return ki < kj
			}, nil)
			path := strings.Join(values, ".")

			if above.lnType == KV_ARR_OPEN {

				if _, ok := mPathIdx[path]; !ok {
					mPathIdx[path] = 0
				}
				// fmt.Printf("\n-----> %v|%v\n\n", path, mPathIdx[path])
				mLvlKey[this.lvl] = fmt.Sprint(mPathIdx[path])
				clrAfter(this.lvl)

			} else if above.line != JUNK {

				mPathIdx[path]++
				// fmt.Printf("\n-----> %v|%v\n\n", path, mPathIdx[path])
				mLvlKey[this.lvl] = fmt.Sprint(mPathIdx[path])
				clrAfter(this.lvl)
			}
		}

		if In(this.lnType, ELEM, ELEM_STR) {

			clrThisAndAfter(this.lvl)

			_, values := MapToKVs(mLvlKey, func(ki, kj int) bool {
				return ki < kj
			}, nil)
			path := strings.Join(values, ".")

			if _, ok := mPathIdx[path]; !ok {
				mPathIdx[path] = 0
			} else {
				mPathIdx[path]++
			}

			// fmt.Printf("\n=====> %v|%v\n\n", path, mPathIdx[path])
			mLvlKey[this.lvl] = fmt.Sprint(mPathIdx[path])

			clrAfter(this.lvl)
		}

		_, values := MapToKVs(mLvlKey, func(ki, kj int) bool {
			return ki < kj
		}, nil)
		suffix := IF(flagCloseObj, "}", "") + IF(flagCloseArr, "]", "") // + IF(flagEmptyObj, "{}", "") + IF(flagEmptyArr, "[]", "")
		return strings.Join(values, ".") + suffix
	}
}

///////////////////////////////////////////////////////////////

func ScanJsonLine(fPathIn, fPathOut string, opt OptLineProc) error {

	if _, err := jt.FmtFileJS(fPathIn); err != nil {
		return err
	}

	SetCurrentKeyLevel := fnSetCurrentKeyLevel()
	I := 0
	mCheck := make(map[string]struct{})

	fd.FileLineScanEx(fPathIn, 1, 1, JUNK, func(line string, cache []string) (bool, string) {

		defer func() { I++ }()

		lnInfo3 := [3]LineInfo{}
		for i, cacheLine := range cache {
			lnType, key, _, _ := simpleFetch(cacheLine)
			lnInfo3[i] = LineInfo{
				line:   cacheLine,
				ln:     strings.TrimSuffix(strings.TrimSpace(cacheLine), ","),
				key:    key,
				lnType: lnType,
				lvl:    strings.Count(strs.HeadBlank(cacheLine), INDENT), // formatted indent is 2 spaces here
			}
		}

		path := SetCurrentKeyLevel(lnInfo3[0], lnInfo3[1], lnInfo3[2])
		// fmt.Printf("%6d: %v\n", I, path)

		if _, ok := mCheck[path]; !ok {
			mCheck[path] = struct{}{}
		} else {
			log.Fatalf("path validation failed: [%d] '%s'", I, path)
		}

		///////////////////////////////////////////////////////////////////////////

		lnType, k, v, comma := simpleFetch(line)
		c := IF(comma, ",", "")
		hb := strs.HeadBlank(line)

		switch lnType {

		case KV:

			if fn := opt.Fn_KV; fn != nil {
				ok, s := fn(I, path, k, v)
				c = IF(len(s) > 0, c, "")
				return ok, hb + s + c
			}
			return true, fmt.Sprintf(`%s"%v": %v%v`, hb, k, v, c)

		case KV_STR:

			if fn := opt.Fn_KV_Str; fn != nil {
				ok, s := fn(I, path, k, v.(string))
				c = IF(len(s) > 0, c, "")
				return ok, hb + s + c
			}
			return true, fmt.Sprintf(`%s"%v": "%v"%v`, hb, k, v, c)

		case KV_OBJ_OPEN:

			if fn := opt.Fn_KV_Obj_Open; fn != nil {
				ok, s := fn(I, path, k, v.(string))
				c = IF(len(s) > 0, c, "")
				return ok, hb + s + c
			}
			return true, fmt.Sprintf(`%s"%v": %v`, hb, k, v)

		case KV_ARR_OPEN:

			if fn := opt.Fn_KV_Arr_Open; fn != nil {
				ok, s := fn(I, path, k, v.(string))
				c = IF(len(s) > 0, c, "")
				return ok, hb + s + c
			}
			return true, fmt.Sprintf(`%s"%v": %v`, hb, k, v)

		case OBJ:

			if fn := opt.Fn_Obj; fn != nil {
				ok, s := fn(I, path, v.(string))
				c = IF(len(s) > 0, c, "")
				return ok, hb + s + c
			}
			return true, hb + v.(string) + c

		case ARR:

			if fn := opt.Fn_Arr; fn != nil {
				ok, s := fn(I, path, v.(string))
				c = IF(len(s) > 0, c, "")
				return ok, hb + s + c
			}
			return true, hb + v.(string) + c

		case ELEM:

			if fn := opt.Fn_Elem; fn != nil {
				ok, s := fn(I, path, v)
				c = IF(len(s) > 0, c, "")
				return ok, hb + s + c
			}
			return true, fmt.Sprintf("%s%v%s", hb, v, c)

		case ELEM_STR:

			if fn := opt.Fn_Elem_Str; fn != nil {
				ok, s := fn(I, path, v.(string))
				c = IF(len(s) > 0, c, "")
				return ok, hb + s + c
			}
			return true, fmt.Sprintf(`%s"%v"%s`, hb, v, c)

		case UNKNOWN:
			panic("Unknown json line")
		}

		return false, ""

	}, fPathOut)

	if _, err := jt.FmtFileJS(fPathOut); err != nil {
		return err
	}

	if len(fPathOut) > 0 {
		data, err := os.ReadFile(fPathOut)
		if err != nil {
			return err
		}
		if len(data) == 0 {
			return fmt.Errorf("FmtFileJS Error After FileLineScanEx")
		}
	}
	return nil
}

///////////////////////////////////////////////////////////////

type OptLineProc struct {
	Fn_KV          func(I int, path, k string, v any) (bool, string)
	Fn_KV_Str      func(I int, path, k, v string) (bool, string)
	Fn_KV_Obj_Open func(I int, path, k, v string) (bool, string)
	Fn_KV_Arr_Open func(I int, path, k, v string) (bool, string)
	Fn_Obj         func(I int, path, v string) (bool, string)
	Fn_Arr         func(I int, path, v string) (bool, string)
	Fn_Elem        func(I int, path string, v any) (bool, string)
	Fn_Elem_Str    func(I int, path, v string) (bool, string)
}

// func fn_kv(I int, path, k string, v any) (bool, string) {
// 	return true, fmt.Sprintf(`"%v": %v`, k, v)
// }

// func fn_kv_str(I int, path, k, v string) (bool, string) {
// 	return true, fmt.Sprintf(`"%v": "%v"`, k, v)
// }

// func fn_kv_obj_open(I int, path, k, v string) (bool, string) {
// 	return true, fmt.Sprintf(`"%v": %v`, k, v)
// }

// func fn_kv_arr_open(I int, path, k, v string) (bool, string) {
// 	return true, fmt.Sprintf(`"%v": %v`, k, v)
// }

// func fn_obj(I int, path, v string) (bool, string) {
// 	return true, v
// }

// func fn_arr(I int, path, v string) (bool, string) {
// 	return true, v
// }

// func fn_elem(I int, path string, v any) (bool, string) {
// 	return true, v.(string)
// }

// func fn_elem_str(I int, path, v string) (bool, string) {
// 	return true, fmt.Sprintf(`"%v"`, v)
// }
