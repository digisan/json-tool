package jsontool

import (
	"encoding/json"
	"strings"

	. "github.com/digisan/go-generics/v2"
)

// IsValid :
func IsValid(bytes []byte) bool {
	var m any
	return json.Unmarshal(bytes, &m) == nil
}

func IsValidStr(str string) bool {
	return IsValid([]byte(str))
}

// Minimize :
func Minimize(str string, check bool) string {

	failOnErrWhen(check && !IsValidStr(str), "%v", fEf("input string is invalid json string"))

	var (
		sb     = &strings.Builder{}
		pc     = byte(0)
		quotes = false
	)

	for i := 0; i < len(str); i++ {
		c := str[i]
		switch {
		case c == '"' && pc != '\\':
			quotes = !quotes
			sb.WriteByte(c)
		case !quotes:
			if NotIn(c, ' ', '\t', '\n', '\r') {
				sb.WriteByte(c)
			}
		case quotes:
			sb.WriteByte(c)
		}
		pc = c
	}

	return sb.String()
}

func TryMinimize(str string) string {
	if !IsValidStr(str) {
		return str
	}
	return Minimize(str, false)
}

// MarshalRemove :
func MarshalRemove(v any, mFieldOldNew map[string]string, rmFields ...string) (bytes []byte, err error) {
	if bytes, err = json.Marshal(v); err != nil {
		return nil, err
	}
	m := make(map[string]any)
	json.Unmarshal(bytes, &m)
	for _, f := range rmFields {
		delete(m, f)
	}
NEXT_NEW:
	for fOld, fNew := range mFieldOldNew {
		for f, v := range m {
			if f == fOld {
				m[fNew] = v
				delete(m, fOld)
				continue NEXT_NEW
			}
		}
	}
	if bytes, err = json.Marshal(m); err != nil {
		return nil, err
	}
	return bytes, nil
}
