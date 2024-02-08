package validator

import (
	. "github.com/digisan/go-generics"
	dt "github.com/digisan/gotk/data-type"
	"github.com/tidwall/gjson"
)

func IsMissing(r gjson.Result) bool {
	return r.Type == gjson.Null && len(r.Raw) == 0
}

func Exists(r gjson.Result) bool {
	return !IsMissing(r)
}

func IsNull(r gjson.Result) bool {
	return r.Type == gjson.Null && r.Raw == "null"
}

func HasValue(r gjson.Result) bool {
	return Exists(r) && !IsNull(r)
}

func HasSomeValue(r gjson.Result) bool {
	return HasValue(r) && !IsEmptyStr(r) && !IsEmptyArr(r) && !IsNestedEmptyArr(r) && !IsAllNullElemArr(r) && !IsAllEmptyStrElemArr(r)
}

func IsBool(r gjson.Result) bool {
	return r.IsBool()
}

func IsNum(r gjson.Result) bool {
	return r.Type == gjson.Number
}

func IsObj(r gjson.Result) bool {
	return r.IsObject()
}

////////////////////////

func IsStr(r gjson.Result) bool {
	return r.Type == gjson.String
}

func IsEmptyStr(r gjson.Result) bool {
	return IsStr(r) && len(r.Str) == 0
}

func IsHTMLStr(r gjson.Result) bool {
	return IsStr(r) && dt.IsHTML([]byte(r.Str))
}

func IsPlainStr(r gjson.Result) bool {
	return IsStr(r) && !dt.IsHTML([]byte(r.Str))
}

////////////////////////

func IsArr(r gjson.Result) bool {
	return r.IsArray()
}

func IsNestedArr(r gjson.Result) bool {
	if IsArr(r) {
		for _, e := range r.Array() {
			if !IsArr(e) {
				return false
			}
		}
		return true
	}
	return false
}

func IsEmptyArr(r gjson.Result) bool {
	return r.IsArray() && len(r.Array()) == 0
}

func IsAllNullElemArr(r gjson.Result) bool {
	if IsArr(r) {
		for _, e := range r.Array() {
			if !IsNull(e) {
				return false
			}
		}
		return true
	}
	return false
}

func IsAllEmptyStrElemArr(r gjson.Result) bool {
	if IsArr(r) {
		for _, e := range r.Array() {
			if !IsEmptyStr(e) {
				return false
			}
		}
		return true
	}
	return false
}

func IsNestedEmptyArr(r gjson.Result) bool {
	if IsArr(r) {
		for _, e := range r.Array() {
			if !IsEmptyArr(e) {
				return false
			}
		}
		return true
	}
	return false
}

func IsStrArr(r gjson.Result) bool {
	if IsArr(r) {
		for _, e := range r.Array() {
			if !IsStr(e) {
				return false
			}
		}
		return true
	}
	return false
}

func IsNestedStrArr(r gjson.Result) bool {
	if IsArr(r) {
		for _, e := range r.Array() {
			if !IsStrArr(e) {
				return false
			}
		}
		return true
	}
	return false
}

func IsPlainStrArr(r gjson.Result) bool {
	if IsStrArr(r) {
		for _, e := range r.Array() {
			if !IsPlainStr(e) {
				return false
			}
		}
		return true
	}
	return false
}

func IsNestedPlainStrArr(r gjson.Result) bool {
	if IsArr(r) {
		for _, e := range r.Array() {
			if !IsPlainStrArr(e) {
				return false
			}
		}
		return true
	}
	return false
}

func IsHTMLStrArr(r gjson.Result) bool {
	if IsStrArr(r) {
		for _, e := range r.Array() {
			if !IsHTMLStr(e) {
				return false
			}
		}
		return true
	}
	return false
}

func IsNestedHTMLStrArr(r gjson.Result) bool {
	if IsArr(r) {
		for _, e := range r.Array() {
			if !IsHTMLStrArr(e) {
				return false
			}
		}
		return true
	}
	return false
}

func IsURLStrArr(r gjson.Result) bool {
	if IsStrArr(r) {
		for _, e := range r.Array() {
			if !IsURL(e.Str) {
				return false
			}
		}
		return true
	}
	return false
}

func IsObjArr(r gjson.Result) bool {
	if IsArr(r) {
		for _, e := range r.Array() {
			if !IsObj(e) {
				return false
			}
		}
		return true
	}
	return false
}

func IsNestedObjArr(r gjson.Result) bool {
	if IsArr(r) {
		for _, e := range r.Array() {
			if !IsObjArr(e) {
				return false
			}
		}
		return true
	}
	return false
}
