package validator

import (
	"strings"

	. "github.com/digisan/go-generics"
	dt "github.com/digisan/gotk/data-type"
	"github.com/tidwall/gjson"
)

func IsMissing(r gjson.Result) bool {
	return r.Type == gjson.Null && len(r.Raw) == 0
}

var NotExist = IsMissing

func Exists(r gjson.Result) bool {
	return !IsMissing(r)
}

func IsNull(r gjson.Result) bool {
	return r.Type == gjson.Null && r.Raw == "null"
}

func IsBool(r gjson.Result) bool {
	return r.IsBool()
}

func IsNum(r gjson.Result) bool {
	return r.Type == gjson.Number
}

func IsStr(r gjson.Result) bool {
	return r.Type == gjson.String
}

func IsObj(r gjson.Result) bool {
	return r.IsObject()
}

func IsArr(r gjson.Result) bool {
	return r.IsArray()
}

func IsEmptyStr(r gjson.Result) bool {
	return IsStr(r) && len(r.Str) == 0
}

func IsEmptyObj(r gjson.Result) bool {
	check := Filter([]byte(r.Raw), func(i int, e byte) bool {
		return NotIn(e, ' ', '\t', '\n')
	})
	return IsObj(r) && string(check) == "{}"
}

func IsEmptyArr(r gjson.Result) bool {
	// return r.IsArray() && r.Raw == "[]"
	return r.IsArray() && len(r.Array()) == 0
}

////////////////////////

func HasNotNullValue(r gjson.Result) bool {
	return Exists(r) && !IsNull(r)
}

func HasSomeValue(r gjson.Result) bool {
	if NotExist(r) {
		return false
	}
	if IsNull(r) {
		return false
	}
	if IsEmptyStr(r) {
		return false
	}
	if IsEmptyObj(r) {
		return false
	}
	if IsEmptyArr(r) {
		return false
	}
	return true
}

func HasMeaningfulValue(r gjson.Result) bool {
	if HasSomeValue(r) {
		return !IsBlankStr(r)
	}
	return false
}

func HasEmptyValue(r gjson.Result) bool {
	return HasNotNullValue(r) && !HasSomeValue(r)
}

////////////////////////

func IsBlankStr(r gjson.Result) bool {
	return IsStr(r) && !IsEmptyStr(r) && len(strings.TrimSpace(r.Str)) == 0
}

func IsHTMLStr(r gjson.Result) bool {
	return IsStr(r) && dt.IsHTML([]byte(r.Str))
}

func IsPlainStr(r gjson.Result) bool {
	return IsStr(r) && !dt.IsHTML([]byte(r.Str))
}

////////////////////////

func IsXArray(r gjson.Result, f func(e gjson.Result) bool) bool {
	if IsArr(r) {
		for _, e := range r.Array() {
			if !f(e) {
				return false
			}
		}
		return len(r.Array()) > 0
	}
	return false
}

func IsNullElemArray(r gjson.Result) bool {
	return IsXArray(r, func(e gjson.Result) bool {
		return IsNull(e)
	})
}

func IsEmptyElemArray(r gjson.Result) bool {
	return IsXArray(r, func(e gjson.Result) bool {
		return HasEmptyValue(e)
	})
}

func IsStrArray(r gjson.Result) bool {
	return IsXArray(r, func(e gjson.Result) bool {
		return IsStr(e)
	})
}

func IsPlainStrArray(r gjson.Result) bool {
	return IsXArray(r, func(e gjson.Result) bool {
		return IsPlainStr(e)
	})
}

func IsHTMLStrArray(r gjson.Result) bool {
	return IsXArray(r, func(e gjson.Result) bool {
		return IsHTMLStr(e)
	})
}

func IsURLStrArray(r gjson.Result) bool {
	return IsXArray(r, func(e gjson.Result) bool {
		return IsURL(e.Str)
	})
}

func IsObjArray(r gjson.Result) bool {
	return IsXArray(r, func(e gjson.Result) bool {
		return IsObj(e)
	})
}

func IsNestedArray(r gjson.Result) bool {
	return IsXArray(r, func(e gjson.Result) bool {
		return IsArr(e)
	})
}

////////////////////////

// func IsAllEmptyStrElemArr(r gjson.Result) bool {
// 	if IsArr(r) {
// 		for _, e := range r.Array() {
// 			if !IsEmptyStr(e) {
// 				return false
// 			}
// 		}
// 		return true
// 	}
// 	return false
// }

// func IsNestedEmptyArr(r gjson.Result) bool {
// 	if IsArr(r) {
// 		for _, e := range r.Array() {
// 			if !IsEmptyArr(e) {
// 				return false
// 			}
// 		}
// 		return true
// 	}
// 	return false
// }

// func IsNestedStrArr(r gjson.Result) bool {
// 	if IsArr(r) {
// 		for _, e := range r.Array() {
// 			if !IsStrArr(e) {
// 				return false
// 			}
// 		}
// 		return true
// 	}
// 	return false
// }

// func IsNestedPlainStrArr(r gjson.Result) bool {
// 	if IsArr(r) {
// 		for _, e := range r.Array() {
// 			if !IsPlainStrArr(e) {
// 				return false
// 			}
// 		}
// 		return true
// 	}
// 	return false
// }

// func IsNestedHTMLStrArr(r gjson.Result) bool {
// 	if IsArr(r) {
// 		for _, e := range r.Array() {
// 			if !IsHTMLStrArr(e) {
// 				return false
// 			}
// 		}
// 		return true
// 	}
// 	return false
// }

// func IsNestedObjArr(r gjson.Result) bool {
// 	if IsArr(r) {
// 		for _, e := range r.Array() {
// 			if !IsObjArr(e) {
// 				return false
// 			}
// 		}
// 		return true
// 	}
// 	return false
// }
